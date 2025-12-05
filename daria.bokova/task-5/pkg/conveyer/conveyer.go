package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type ChanMap map[string]chan string

type Conveyer struct {
	mu          sync.RWMutex
	channels    ChanMap
	decorators  []DecoratorConfig
	multiplexer *MultiplexerConfig
	separator   *SeparatorConfig
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	size        int
}

type DecoratorConfig struct {
	fn     func(context.Context, chan string, chan string) error
	input  string
	output string
}

type MultiplexerConfig struct {
	fn     func(context.Context, []chan string, chan string) error
	inputs []string
	output string
}

type SeparatorConfig struct {
	fn      func(context.Context, chan string, []chan string) error
	input   string
	outputs []string
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(ChanMap),
		size:     size,
	}
}

func (c *Conveyer) getOrCreateChan(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[id]; exists {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[id] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.decorators = append(c.decorators, DecoratorConfig{fn, input, output})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.multiplexer = &MultiplexerConfig{fn, inputs, output}
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.separator = &SeparatorConfig{fn, input, outputs}
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.startWorkers(ctx)

	c.wg.Wait()

	c.cleanup()

	return nil
}

func (c *Conveyer) startWorkers(ctx context.Context) {
	for _, d := range c.decorators {
		inputCh := c.getOrCreateChan(d.input)
		outputCh := c.getOrCreateChan(d.output)
		c.wg.Add(1)
		go func(dc DecoratorConfig) {
			defer c.wg.Done()
			if err := dc.fn(ctx, inputCh, outputCh); err != nil {
				// Обработка ошибки - просто логируем
				fmt.Printf("Decorator error: %v\n", err)
			}
		}(d)
	}

	if c.multiplexer != nil {
		var inputs []chan string
		for _, in := range c.multiplexer.inputs {
			inputs = append(inputs, c.getOrCreateChan(in))
		}
		outputCh := c.getOrCreateChan(c.multiplexer.output)
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			if err := c.multiplexer.fn(ctx, inputs, outputCh); err != nil {
				fmt.Printf("Multiplexer error: %v\n", err)
			}
		}()
	}

	if c.separator != nil {
		inputCh := c.getOrCreateChan(c.separator.input)
		var outputs []chan string
		for _, out := range c.separator.outputs {
			outputs = append(outputs, c.getOrCreateChan(out))
		}
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			if err := c.separator.fn(ctx, inputCh, outputs); err != nil {
				fmt.Printf("Separator error: %v\n", err)
			}
		}()
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	ch, exists := c.channels[input]
	c.mu.RUnlock()
	if !exists {
		return errors.New("chan not found")
	}

	select {
	case ch <- data:
		return nil
	default:
		return errors.New("channel is full")
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[output]
	c.mu.RUnlock()
	if !exists {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
}

func (c *Conveyer) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		select {
		case _, ok := <-ch:
			if !ok {
				continue
			}
		default:
			close(ch)
		}
	}
	c.channels = make(ChanMap)
}

func (c *Conveyer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}
