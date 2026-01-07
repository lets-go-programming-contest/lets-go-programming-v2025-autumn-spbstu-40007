package conveyer

import (
	"context"
	"fmt"
	"sync"
)

type Conveyer struct {
	size      int
	channels  map[string]chan string
	handlers  []func(ctx context.Context) error
	mu        sync.RWMutex
	isRunning bool
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
	}
}

func (c *Conveyer) getChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input, output string,
) {
	in := c.getChan(input)
	out := c.getChan(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string, output string,
) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = c.getChan(name)
	}
	out := c.getChan(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChans, out)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string, outputs []string,
) {
	in := c.getChan(input)
	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = c.getChan(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	if c.isRunning {
		return fmt.Errorf("already running")
	}
	c.isRunning = true

	errChan := make(chan error, len(c.handlers))

	var wg sync.WaitGroup
	for _, handler := range c.handlers {
		wg.Add(1)
		go func(h func(context.Context) error) {
			defer wg.Done()
			if err := h(ctx); err != nil {
				errChan <- err
			}
		}(handler)
	}

	go func() {
		wg.Wait()
		close(errChan)

		c.mu.Lock()
		for _, ch := range c.channels {
			close(ch)
		}
		c.mu.Unlock()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err, ok := <-errChan:
		if ok && err != nil {
			return err
		}
		return nil
	}
}

func (c *Conveyer) Send(input, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return fmt.Errorf("chan not found")
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("chan not found")
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
