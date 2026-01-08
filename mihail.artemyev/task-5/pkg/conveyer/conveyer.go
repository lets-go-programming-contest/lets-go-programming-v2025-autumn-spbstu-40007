package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type Conveyer struct {
	size     int
	mu       sync.RWMutex
	channels map[string]chan string
	handlers []func(context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		mu:       sync.RWMutex{},
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
	}
}

func (c *Conveyer) register(name string) chan string {
	c.mu.RLock()
	ch, exists := c.channels[name]
	c.mu.RUnlock()

	if !exists {
		c.mu.Lock()
		if chLocal, ok := c.channels[name]; ok {
			ch = chLocal
		} else {
			ch = make(chan string, c.size)
			c.channels[name] = ch
		}
		c.mu.Unlock()
	}
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
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
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChans := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		inputChans = append(inputChans, c.register(name))
	}

	out := c.register(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChans, out)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	in := c.register(input)

	outputChans := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outputChans = append(outputChans, c.register(name))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outputChans)
	})
}

func closeChannelSafe(ch chan string) {
	defer func() {
		recover()
	}()
	close(ch)
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		closeChannelSafe(ch)
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.RLock()
	handlers := make([]func(context.Context) error, len(c.handlers))
	copy(handlers, c.handlers)
	c.mu.RUnlock()

	group, ctxWithCancel := errgroup.WithContext(ctx)

	for _, h := range handlers {
		handler := h
		group.Go(func() error {
			return handler(ctxWithCancel)
		})
	}

	err := group.Wait()

	c.closeAllChannels()

	if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("conveyer failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(inputName string, data string) error {
	channel, exists := func() (chan string, bool) {
		c.mu.RLock()
		defer c.mu.RUnlock()
		channel, ok := c.channels[inputName]
		return channel, ok
	}()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data
	return nil
}

func (c *Conveyer) Recv(outputName string) (string, error) {
	channel, exists := func() (chan string, bool) {
		c.mu.RLock()
		defer c.mu.RUnlock()
		channel, ok := c.channels[outputName]
		return channel, ok
	}()

	if !exists {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return value, nil
}
