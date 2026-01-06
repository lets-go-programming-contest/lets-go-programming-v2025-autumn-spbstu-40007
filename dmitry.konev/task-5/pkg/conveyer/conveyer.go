package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer interface {
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

type handlerRunner func(context.Context) error

type conveyer struct {
	size     int
	mu       sync.RWMutex
	chans    map[string]chan string
	handlers []handlerRunner
}

func New(size int) Conveyer {
	return &conveyer{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (c *conveyer) getOrCreate(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.chans[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.chans[name] = ch
	return ch
}

func (c *conveyer) RegisterDecorator(
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

func (c *conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	ins := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		ins = append(ins, c.getOrCreate(name))
	}
	out := c.getOrCreate(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (c *conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	in := c.getOrCreate(input)

	outs := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outs = append(outs, c.getOrCreate(name))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outs)
	})
}

func (c *conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, h := range c.handlers {
		wg.Add(1)
		go func(h handlerRunner) {
			defer wg.Done()
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
	case err := <-errCh:
		wg.Wait()
		c.closeAll()
		return err
	}

	wg.Wait()
	c.closeAll()
	return nil
}

func (c *conveyer) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.chans {
		close(ch)
	}
}

func (c *conveyer) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.chans[input]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.chans[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return val, nil
}
