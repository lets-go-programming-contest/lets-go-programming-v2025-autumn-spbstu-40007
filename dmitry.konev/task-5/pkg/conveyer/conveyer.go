package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
)

type decoratorFunc func(context.Context, chan string, chan string) error
type multiplexerFunc func(context.Context, []chan string, chan string) error
type separatorFunc func(context.Context, chan string, []chan string) error

type handler struct {
	run func(context.Context) error
}

type Conveyor interface {
	RegisterDecorator(fn decoratorFunc, input string, output string)
	RegisterMultiplexer(fn multiplexerFunc, inputs []string, output string)
	RegisterSeparator(fn separatorFunc, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyor struct {
	size     int
	channels map[string]chan string
	handlers []handler
	mu       sync.Mutex
}

func New(size int) *conveyor {
	return &conveyor{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]handler, 0),
	}
}

func (c *conveyor) getOrCreate(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.channels[name]
	if ok {
		return ch
	}

	ch = make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyor) get(name string) (chan string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.channels[name]
	return ch, ok
}

func (c *conveyor) RegisterDecorator(fn decoratorFunc, input, output string) {
	inCh := c.getOrCreate(input)
	outCh := c.getOrCreate(output)
	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inCh, outCh)
		},
	})
}

func (c *conveyor) RegisterMultiplexer(fn multiplexerFunc, inputs []string, output string) {
	inChs := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChs[i] = c.getOrCreate(name)
	}
	outCh := c.getOrCreate(output)
	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inChs, outCh)
		},
	})
}

func (c *conveyor) RegisterSeparator(fn separatorFunc, input string, outputs []string) {
	inCh := c.getOrCreate(input)
	outChs := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChs[i] = c.getOrCreate(name)
	}
	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inCh, outChs)
		},
	})
}

func (c *conveyor) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		h := h
		group.Go(func() error {
			return h.run(ctx)
		})
	}

	err := group.Wait()

	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mu.Unlock()

	return err
}

func (c *conveyor) Send(input string, data string) error {
	ch, ok := c.get(input)
	if !ok {
		return ErrChanNotFound
	}

	select {
	case ch <- data:
		return nil
	default:
		ch <- data
		return nil
	}
}

func (c *conveyor) Recv(output string) (string, error) {
	ch, ok := c.get(output)
	if !ok {
		return "", ErrChanNotFound
	}

	v, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return v, nil
}