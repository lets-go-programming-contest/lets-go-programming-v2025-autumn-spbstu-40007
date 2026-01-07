package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrUndefined    = errors.New("undefined")
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
	}
}

func (c *conveyor) getOrCreate(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
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
	in := c.getOrCreate(input)
	out := c.getOrCreate(output)
	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, in, out)
		},
	})
}

func (c *conveyor) RegisterMultiplexer(fn multiplexerFunc, inputs []string, output string) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = c.getOrCreate(name)
	}
	out := c.getOrCreate(output)
	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inChans, out)
		},
	})
}

func (c *conveyor) RegisterSeparator(fn separatorFunc, input string, outputs []string) {
	in := c.getOrCreate(input)
	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = c.getOrCreate(name)
	}
	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, in, outChans)
		},
	})
}

func (c *conveyor) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	for _, h := range c.handlers {
		hCopy := h
		group.Go(func() error {
			return hCopy.run(ctx)
		})
	}
	return group.Wait()
}

func (c *conveyor) Send(input, data string) error {
	ch, ok := c.get(input)
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *conveyor) Recv(output string) (string, error) {
	ch, ok := c.get(output)
	if !ok {
		return "", ErrChanNotFound
	}
	val, ok := <-ch
	if !ok {
		return "", ErrUndefined
	}
	return val, nil
}