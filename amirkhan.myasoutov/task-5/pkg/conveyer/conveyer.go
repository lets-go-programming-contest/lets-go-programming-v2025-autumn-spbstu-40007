package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound   = errors.New("chan not found")
	ErrInvalidHandler = errors.New("invalid handler function signature or nil")
)

type DecoratorHandler func(ctx context.Context, input chan string, output chan string) error
type MultiplexerHandler func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorHandler func(ctx context.Context, input chan string, outputs []chan string) error

type HandlerConfig struct {
	HandlerFn interface{}
	InputIDs  []string
	OutputIDs []string
}

type Pipeline interface {
	RegisterDecorator(fn DecoratorHandler, inputID string, outputID string)
	RegisterMultiplexer(fn MultiplexerHandler, inputsIDs []string, outputID string)
	RegisterSeparator(fn SeparatorHandler, inputID string, outputsIDs []string)
	Run(ctx context.Context) error
	Send(inputID string, data string) error
	Recv(outputID string) (string, error)
}

type conveyerImpl struct {
	channelSize    int
	channels       map[string]chan string
	handlerConfigs []HandlerConfig
	mu             sync.Mutex
	wg             sync.WaitGroup
}

func New(size int) Pipeline {
	return &conveyerImpl{
		channelSize:    size,
		channels:       make(map[string]chan string),
		handlerConfigs: make([]HandlerConfig, 0),
	}
}

func (c *conveyerImpl) getChan(id string) chan string {

	if ch, ok := c.channels[id]; ok {
		return ch
	}

	newChan := make(chan string, c.channelSize)
	c.channels[id] = newChan
	return newChan
}

func (c *conveyerImpl) getChannelsByIds(ids []string) []chan string {
	channels := make([]chan string, len(ids))
	for i, id := range ids {
		channels[i] = c.getChan(id)
	}
	return channels
}

func (c *conveyerImpl) RegisterDecorator(fn DecoratorHandler, inputID string, outputID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlerConfigs = append(c.handlerConfigs, HandlerConfig{
		HandlerFn: fn,
		InputIDs:  []string{inputID},
		OutputIDs: []string{outputID},
	})
	c.getChan(inputID)
	c.getChan(outputID)
}

func (c *conveyerImpl) RegisterMultiplexer(fn MultiplexerHandler, inputsIDs []string, outputID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlerConfigs = append(c.handlerConfigs, HandlerConfig{
		HandlerFn: fn,
		InputIDs:  inputsIDs,
		OutputIDs: []string{outputID},
	})
	for _, id := range inputsIDs {
		c.getChan(id)
	}
	c.getChan(outputID)
}

func (c *conveyerImpl) RegisterSeparator(fn SeparatorHandler, inputID string, outputsIDs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlerConfigs = append(c.handlerConfigs, HandlerConfig{
		HandlerFn: fn,
		InputIDs:  []string{inputID},
		OutputIDs: outputsIDs,
	})
	c.getChan(inputID)
	for _, id := range outputsIDs {
		c.getChan(id)
	}
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	errChan := make(chan error, len(c.handlerConfigs))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, config := range c.handlerConfigs {
		c.wg.Add(1)

		inputChans := c.getChannelsByIds(config.InputIDs)
		outputChans := c.getChannelsByIds(config.OutputIDs)

		go func(cfg HandlerConfig, ins []chan string, outs []chan string) {
			defer c.wg.Done()
			var err error

			switch fn := cfg.HandlerFn.(type) {
			case DecoratorHandler:
				if len(ins) != 1 || len(outs) != 1 {
					err = errors.New("decorator handler requires 1 input and 1 output")
				} else {
					err = fn(ctx, ins[0], outs[0])
				}
			case MultiplexerHandler:
				if len(outs) != 1 {
					err = errors.New("multiplexer handler requires 1 output")
				} else {
					err = fn(ctx, ins, outs[0])
				}
			case SeparatorHandler:
				if len(ins) != 1 || len(outs) == 0 {
					err = errors.New("separator handler requires 1 input and >=1 output")
				} else {
					err = fn(ctx, ins[0], outs)
				}
			default:
				err = ErrInvalidHandler
			}

			if err != nil && err != context.Canceled {
				errChan <- fmt.Errorf("handler failed: %w", err)
			}
		}(config, inputChans, outputChans)
	}

	select {
	case <-ctx.Done():
		c.wg.Wait()
		return ctx.Err()
	case err := <-errChan:
		cancel()
		c.wg.Wait()
		return err
	}
}

func (c *conveyerImpl) Send(inputID string, data string) error {
	c.mu.Lock()
	ch, ok := c.channels[inputID]
	c.mu.Unlock()

	if !ok {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(outputID string) (string, error) {
	c.mu.Lock()
	ch, ok := c.channels[outputID]
	c.mu.Unlock()

	if !ok {
		return "", ErrChanNotFound
	}

	data, open := <-ch
	if !open {
		return "undefined", nil
	}
	return data, nil
}
