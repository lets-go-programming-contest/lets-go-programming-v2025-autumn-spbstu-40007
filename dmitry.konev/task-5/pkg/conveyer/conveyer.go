package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyor interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		handler func(ctx context.Context, input chan string, output chan string),
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		handler func(ctx context.Context, inputs []chan string, output chan string),
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		handler func(ctx context.Context, input chan string, outputs []chan string),
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyorImpl struct {
	chans   map[string]chan string
	size    int
	workers []func(ctx context.Context) error
}

func New(size int) *conveyorImpl {
	return &conveyorImpl{
		chans:   make(map[string]chan string),
		size:    size,
		workers: make([]func(ctx context.Context) error, 0),
	}
}

func (c *conveyorImpl) getOrCreateChan(name string) chan string {
	if ch, ok := c.chans[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.chans[name] = ch
	return ch
}

func (c *conveyorImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	_ func(ctx context.Context, input chan string, output chan string),
	input string,
	output string,
) {
	in := c.getOrCreateChan(input)
	out := c.getOrCreateChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *conveyorImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	_ func(ctx context.Context, inputs []chan string, output chan string),
	inputs []string,
	output string,
) {
	in := make([]chan string, len(inputs))
	for i, name := range inputs {
		in[i] = c.getOrCreateChan(name)
	}
	out := c.getOrCreateChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *conveyorImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	_ func(ctx context.Context, input chan string, outputs []chan string),
	input string,
	outputs []string,
) {
	in := c.getOrCreateChan(input)
	out := make([]chan string, len(outputs))
	for i, name := range outputs {
		out[i] = c.getOrCreateChan(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *conveyorImpl) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.workers))

	for _, w := range c.workers {
		worker := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := worker(ctx); err != nil && !errors.Is(err, context.Canceled) {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	wg.Wait()
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (c *conveyorImpl) Send(input string, data string) error {
	ch, ok := c.chans[input]
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *conveyorImpl) Recv(output string) (string, error) {
	ch, ok := c.chans[output]
	if !ok {
		return "", ErrChanNotFound
	}
	val, ok := <-ch
	if !ok {
		return "", nil
	}
	return val, nil
}
