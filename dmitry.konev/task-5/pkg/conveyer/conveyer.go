package conveyer

import (
	"context"
	"errors"
	"sync"
)

const (
	ErrChanNotFound = "chan not found"
	Undefined       = "undefined"
)

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	mu       sync.RWMutex

	workers []func(ctx context.Context) error
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *conveyerImpl) getOrCreate(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[id]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[id] = ch
	return ch
}

func (c *conveyerImpl) get(id string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.channels[id]
	return ch, ok
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	in := c.getOrCreate(input)
	out := c.getOrCreate(output)
	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	out := c.getOrCreate(output)
	ins := make([]chan string, 0, len(inputs))

	for _, id := range inputs {
		ins = append(ins, c.getOrCreate(id))
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	in := c.getOrCreate(input)
	outs := make([]chan string, 0, len(outputs))

	for _, id := range outputs {
		outs = append(outs, c.getOrCreate(id))
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, outs)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, len(c.workers))
	var wg sync.WaitGroup

	for _, w := range c.workers {
		wg.Add(1)
		go func(fn func(context.Context) error) {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				errCh <- err
			}
		}(w)
	}

	select {
	case <-ctx.Done():
	case err := <-errCh:
		cancel()
		wg.Wait()
		c.closeAll()
		return err
	}

	wg.Wait()
	c.closeAll()
	return nil
}

func (c *conveyerImpl) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, ok := c.get(input)
	if !ok {
		return errors.New(ErrChanNotFound)
	}
	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, ok := c.get(output)
	if !ok {
		return "", errors.New(ErrChanNotFound)
	}

	val, ok := <-ch
	if !ok {
		return Undefined, nil
	}
	return val, nil
}
