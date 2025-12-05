package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrClosed       = errors.New("channel closed")
)

type handlerFn func(ctx context.Context) error

const undefined = "undefined"

type Conveyer struct {
	size int

	mu       sync.RWMutex
	channels map[string]chan string
	handlers []handlerFn
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]handlerFn, 0),
	}
}

func (c *Conveyer) register(name string) chan string {
	if _, exists := c.channels[name]; !exists {
		c.channels[name] = make(chan string, c.size)
	}
	return c.channels[name]
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	in := c.register(input)
	out := c.register(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var inputChans []chan string
	for _, name := range inputs {
		inputChans = append(inputChans, c.register(name))
	}
	out := c.register(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChans, out)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	in := c.register(input)
	var outputChans []chan string
	for _, name := range outputs {
		outputChans = append(outputChans, c.register(name))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outputChans)
	})
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()
			close(ch)
		}()
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	errGr, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		handler := handler
		errGr.Go(func() error {
			return handler(ctx)
		})
	}

	err := errGr.Wait()

	c.closeAllChannels()

	if err != nil && err != context.Canceled {
		return fmt.Errorf("conveyer failed: %w", err)
	}

	if err == context.Canceled && ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	channel, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	select {
	case channel <- data:
		return nil
	default:
		return fmt.Errorf("send to channel %s failed: channel blocked", input)
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return undefined, ErrChanNotFound
	}

	val, ok := <-channel
	if !ok {
		return undefined, nil
	}
	return val, nil
}
