package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound   = errors.New("chan not found")
	ErrUndefined      = errors.New("undefined")
	ErrClosed         = errors.New("closed")
	ErrAlreadyStarted = errors.New("already started")
)

type decoratorDesc struct {
	fn         func(ctx context.Context, input chan string, output chan string, errCh chan error)
	inputName  string
	outputName string
}

type separatorDesc struct {
	fn          func(ctx context.Context, input chan string, outputs []chan string, errCh chan error)
	inputName   string
	outputNames []string
}

type multiplexerDesc struct {
	fn         func(ctx context.Context, inputs []chan string, output chan string, errCh chan error)
	inputNames []string
	outputName string
}

type Conveyer interface {
	RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string, errCh chan error),
		input string, output string)
	RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string, errCh chan error),
		inputs []string, output string)
	RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string, errCh chan error),
		input string, outputs []string)

	Send(input string, data string) error
	Recv(output string) (string, error)
	Run(ctx context.Context) error
}

type conveyorImpl struct {
	size    int
	mu      sync.RWMutex
	chans   map[string]chan string
	decs    []decoratorDesc
	seps    []separatorDesc
	muxes   []multiplexerDesc
	errCh   chan error
	started bool
}

func New(size int) *conveyorImpl {
	return &conveyorImpl{
		size:    size,
		chans:   make(map[string]chan string),
		decs:    make([]decoratorDesc, 0),
		seps:    make([]separatorDesc, 0),
		muxes:   make([]multiplexerDesc, 0),
		errCh:   make(chan error, size),
		started: false,
	}
}

func (c *conveyorImpl) getOrCreate(name string) chan string {
	channel, ok := c.chans[name]
	if !ok {
		channel = make(chan string, c.size)
		c.chans[name] = channel
	}

	return channel
}

func (c *conveyorImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string, errCh chan error),
	input string, output string,
) {
	c.decs = append(c.decs, decoratorDesc{
		fn:         fn,
		inputName:  input,
		outputName: output,
	})
}

func (c *conveyorImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string, errCh chan error),
	input string, outputs []string,
) {
	c.seps = append(c.seps, separatorDesc{
		fn:          fn,
		inputName:   input,
		outputNames: outputs,
	})
}

func (c *conveyorImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string, errCh chan error),
	inputs []string, output string,
) {
	c.muxes = append(c.muxes, multiplexerDesc{
		fn:         fn,
		inputNames: inputs,
		outputName: output,
	})
}

func (c *conveyorImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()

		return ErrAlreadyStarted
	}
	c.started = true
	c.mu.Unlock()

	var waitGroup sync.WaitGroup

	for _, decorator := range c.decs {
		inputCh := c.getOrCreate(decorator.inputName)
		outputCh := c.getOrCreate(decorator.outputName)

		waitGroup.Add(1)

		go func(desc decoratorDesc, input, output chan string) {
			defer waitGroup.Done()
			desc.fn(ctx, input, output, c.errCh)
		}(decorator, inputCh, outputCh)
	}

	for _, separator := range c.seps {
		inputCh := c.getOrCreate(separator.inputName)
		outs := make([]chan string, 0, len(separator.outputNames))

		for _, name := range separator.outputNames {
			outs = append(outs, c.getOrCreate(name))
		}

		waitGroup.Add(1)

		go func(desc separatorDesc, input chan string, outputs []chan string) {
			defer waitGroup.Done()
			desc.fn(ctx, input, outputs, c.errCh)
		}(separator, inputCh, outs)
	}

	for _, multiplexer := range c.muxes {
		inputs := make([]chan string, 0, len(multiplexer.inputNames))
		for _, name := range multiplexer.inputNames {
			inputs = append(inputs, c.getOrCreate(name))
		}
		outputCh := c.getOrCreate(multiplexer.outputName)

		waitGroup.Add(1)

		go func(desc multiplexerDesc, ins []chan string, output chan string) {
			defer waitGroup.Done()
			desc.fn(ctx, ins, output, c.errCh)
		}(multiplexer, inputs, outputCh)
	}

	go func() {
		waitGroup.Wait()
		close(c.errCh)
	}()

	return nil
}

func (c *conveyorImpl) Send(input string, data string) error {
	c.mu.RLock()
	channel, ok := c.chans[input]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *conveyorImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, ok := c.chans[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return "", ErrClosed
	}

	return value, nil
}
