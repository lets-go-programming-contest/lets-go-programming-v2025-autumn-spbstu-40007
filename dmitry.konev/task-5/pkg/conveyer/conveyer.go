package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		handlerFunc func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)

	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)

	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
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
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputch := c.getOrCreateChan(input)
	outputch := c.getOrCreateChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		if err := fn(ctx, inputch, outputch); err != nil {
			return err
		}
		return handlerFunc(ctx, inputch, outputch)
	})
}

func (c *conveyorImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputsch := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputsch[i] = c.getOrCreateChan(name)
	}

	outputch := c.getOrCreateChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		if err := fn(ctx, inputsch, outputch); err != nil {
			return err
		}
		return handlerFunc(ctx, inputsch, outputch)
	})
}

func (c *conveyorImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputch := c.getOrCreateChan(input)
	outputsCh := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputsCh[i] = c.getOrCreateChan(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		if err := fn(ctx, inputch, outputsCh); err != nil {
			return err
		}
		return handlerFunc(ctx, inputch, outputsCh)
	})
}

func (c *conveyorImpl) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	errch := make(chan error, len(c.workers))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defer func() {
		for _, ch := range c.chans {
			close(ch)
		}
	}()

	for _, worker := range c.workers {
		wg.Add(1)
		w := worker

		go func() {
			defer wg.Done()
			if err := w(ctx); err != nil {
				if !errors.Is(err, context.Canceled) {
					select {
					case errch <- err:
						cancel()
					default:
					}
				}
			}
		}()
	}

	select {
	case <-ctx.Done():
		wg.Wait()
		return nil
	case err := <-errch:
		wg.Wait()
		return err
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

	val, isOpen := <-ch
	if !isOpen {
		return "undefined", nil
	}

	return val, nil
}
