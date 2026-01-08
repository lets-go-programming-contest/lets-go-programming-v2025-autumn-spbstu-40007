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
		handler func(ctx context.Context, input chan string, output chan string),
		source string,
		sink string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		handler func(ctx context.Context, inputs []chan string, output chan string),
		sources []string,
		sink string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		handler func(ctx context.Context, input chan string, outputs []chan string),
		source string,
		sinks []string,
	)
	Run(ctx context.Context) error
	Send(source string, val string) error
	Recv(sink string) (string, error)
}

type conveyorImpl struct {
	channels map[string]chan string
	buffer   int
	tasks    []func(ctx context.Context) error
}

func New(buffer int) *conveyorImpl {
	return &conveyorImpl{
		channels: make(map[string]chan string),
		buffer:   buffer,
		tasks:    make([]func(ctx context.Context) error, 0),
	}
}

func (c *conveyorImpl) getChannel(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.buffer)
	c.channels[name] = ch
	return ch
}

func (c *conveyorImpl) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.tasks))

	for _, task := range c.tasks {
		t := task
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := t(ctx); err != nil && !errors.Is(err, context.Canceled) {
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

func (c *conveyorImpl) Send(source string, val string) error {
	ch, ok := c.channels[source]
	if !ok {
		return ErrChanNotFound
	}
	ch <- val
	return nil
}

func (c *conveyorImpl) Recv(sink string) (string, error) {
	ch, ok := c.channels[sink]
	if !ok {
		return "", ErrChanNotFound
	}
	v, open := <-ch
	if !open {
		return "", nil
	}
	return v, nil
}
