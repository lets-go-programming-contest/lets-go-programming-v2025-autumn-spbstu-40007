package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChannelNotFound = errors.New("chan not found")

type Conveyer interface {
	RegisterDecorator(
		handlerFunc func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type implementation struct {
	channels map[string]chan string
	size     int
	workers  []func(ctx context.Context) error
}

func New(size int) *implementation {
	return &implementation{
		channels: make(map[string]chan string),
		size:     size,
		workers:  make([]func(ctx context.Context) error, 0),
	}
}

func (c *implementation) getOrCreateChannel(channelID string) chan string {
	if ch, exists := c.channels[channelID]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[channelID] = ch

	return ch
}

func (c *implementation) RegisterDecorator(
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inCh := c.getOrCreateChannel(input)
	outCh := c.getOrCreateChannel(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outCh)
	})
}

func (c *implementation) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inChs := make([]chan string, len(inputs))
	for i, id := range inputs {
		inChs[i] = c.getOrCreateChannel(id)
	}

	outCh := c.getOrCreateChannel(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return handlerFunc(ctx, inChs, outCh)
	})
}

func (c *implementation) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.getOrCreateChannel(input)
	outChs := make([]chan string, len(outputs))

	for i, id := range outputs {
		outChs[i] = c.getOrCreateChannel(id)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outChs)
	})
}

func (c *implementation) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup

	errCh := make(chan error, len(c.workers))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defer func() {
		for _, ch := range c.channels {
			close(ch)
		}
	}()

	for _, worker := range c.workers {
		waitGroup.Add(1)

		w := worker

		go func() {
			defer waitGroup.Done()

			if err := w(ctx); err != nil {
				if !errors.Is(err, context.Canceled) {
					select {
					case errCh <- err:
						cancel()
					default:
					}
				}
			}
		}()
	}

	select {
	case <-ctx.Done():
		waitGroup.Wait()

		return nil
	case err := <-errCh:
		waitGroup.Wait()

		return err
	}
}

func (c *implementation) Send(input string, data string) error {
	ch, exists := c.channels[input]
	if !exists {
		return ErrChannelNotFound
	}

	ch <- data

	return nil
}

func (c *implementation) Recv(output string) (string, error) {
	ch, exists := c.channels[output]
	if !exists {
		return "", ErrChannelNotFound
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return val, nil
}
