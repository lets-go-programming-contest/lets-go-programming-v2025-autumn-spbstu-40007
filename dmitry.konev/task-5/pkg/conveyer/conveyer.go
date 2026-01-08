package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer interface {
	RegisterDecorator(
		workerFunc func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)

	RegisterMultiplexer(
		workerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)

	RegisterSeparator(
		workerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
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
	workerFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputch := c.getOrCreateChan(input)
	outputch := c.getOrCreateChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return workerFunc(ctx, inputch, outputch)
	})
}

func (c *conveyorImpl) RegisterMultiplexer(
	workerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputsch := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputsch[i] = c.getOrCreateChan(name)
	}

	outputch := c.getOrCreateChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return workerFunc(ctx, inputsch, outputch)
	})
}

func (c *conveyorImpl) RegisterSeparator(
	workerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputch := c.getOrCreateChan(input)
	outputsCh := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputsCh[i] = c.getOrCreateChan(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return workerFunc(ctx, inputch, outputsCh)
	})
}

func (c *conveyorImpl) Run(ctx context.Context) error {
	var workGroup sync.WaitGroup

	errCh := make(chan error, len(c.workers))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defer func() {
		for _, ch := range c.chans {
			close(ch)
		}
	}()

	for _, worker := range c.workers {
		workGroup.Add(1)
		w := worker

		go func() {
			defer workGroup.Done()
			err := w(ctx)

			if err == nil || errors.Is(err, context.Canceled) {
				return
			}

			select {
			case errCh <- err:
				cancel()
			default:
			}
		}()
	}

	workGroup.Wait()

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

	val, isOpen := <-ch
	if !isOpen {
		return "undefined", nil
	}

	return val, nil
}