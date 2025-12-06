//nolint:varnamelen
package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type handlersFn func(ctx context.Context) error

const undefined = "undefined"

type conveyer struct {
	size int

	mu       sync.RWMutex
	channels map[string]chan string
	handlers []handlersFn
}

func New(size int) *conveyer {
	return &conveyer{ //nolint:exhaustruct
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]handlersFn, 0),
	}
}

func (c *conveyer) register(name string) chan string {
	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *conveyer) RegisterDecorator(
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

func (c *conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
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

func (c *conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
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
		_ = recover()
	}()

	close(ch)
}

func (c *conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		closeChannelSafe(ch)
	}
}

func (c *conveyer) Run(ctx context.Context) error {
	c.mu.RLock()
	handlers := make([]handlersFn, len(c.handlers))
	copy(handlers, c.handlers)
	c.mu.RUnlock()

	group, ctxWithCancel := errgroup.WithContext(ctx)

	for _, handlers := range handlers {
		h := handlers

		group.Go(func() error {
			return h(ctxWithCancel)
		})
	}

	err := group.Wait()

	c.closeAllChannels()

	if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("conveyer failed: %w", err)
	}

	return nil
}

func (c *conveyer) Send(input string, data string) error {
	c.mu.RLock()
	channel, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return value, nil
}
