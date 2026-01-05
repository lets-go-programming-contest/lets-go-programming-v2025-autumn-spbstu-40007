package conveyer

import (
	"context"
	"errors"
	"sync"
)

type conveyer interface {
	RegisterDecorator(
		fn func(context.Context, chan string, chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(context.Context, []chan string, chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(context.Context, chan string, []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type handlerRunner func(ctx context.Context) error

type Conveyer struct {
	size     int
	chans    map[string]chan string
	mu       sync.RWMutex
	handlers []handlerRunner
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (c *Conveyer) getOrCreate(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.chans[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.chans[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	in := c.getOrCreate(input)
	out := c.getOrCreate(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	var ins []chan string
	for _, name := range inputs {
		ins = append(ins, c.getOrCreate(name))
	}
	out := c.getOrCreate(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	in := c.getOrCreate(input)
	var outs []chan string
	for _, name := range outputs {
		outs = append(outs, c.getOrCreate(name))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outs)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	for _, h := range c.handlers {
		go func(h handlerRunner) {
			if err := h(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				cancel()
			}
		}(h)
	}

	select {
	case <-ctx.Done():
		c.closeAll()
		return nil
	case err := <-errCh:
		c.closeAll()
		return err
	}
}

func (c *Conveyer) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.chans {
		close(ch)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.chans[input]
	c.mu.RUnlock()

	if !ok {
		return errors.New("chan not found")
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.chans[output]
	c.mu.RUnlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return val, nil
}