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

	channel, exists := c.channels[name]
	if exists {
		return channel
	}

	channel = make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *conveyor) get(name string) (chan string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, exists := c.channels[name]
	return channel, exists
}

func (c *conveyor) RegisterDecorator(
	fn decoratorFunc,
	input string,
	output string,
) {
	inputCh := c.getOrCreate(input)
	outputCh := c.getOrCreate(output)

	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inputCh, outputCh)
		},
	})
}

func (c *conveyor) RegisterMultiplexer(
	fn multiplexerFunc,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		inputChans = append(inputChans, c.getOrCreate(name))
	}

	outputCh := c.getOrCreate(output)

	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inputChans, outputCh)
		},
	})
}

func (c *conveyor) RegisterSeparator(
	fn separatorFunc,
	input string,
	outputs []string,
) {
	inputCh := c.getOrCreate(input)

	outputChans := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outputChans = append(outputChans, c.getOrCreate(name))
	}

	c.handlers = append(c.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inputCh, outputChans)
		},
	})
}

func (c *conveyor) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		handlerRun := h.run

		group.Go(func() error {
			return handlerRun(ctx)
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
	channel, exists := c.get(input)
	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *conveyor) Recv(output string) (string, error) {
	channel, exists := c.get(output)
	if !exists {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return "", ErrUndefined
	}

	return value, nil
}