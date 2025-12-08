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
	fn         func(ctx context.Context, input chan string, output chan string) error
	inputName  string
	outputName string
}

type separatorDesc struct {
	fn          func(ctx context.Context, input chan string, outputs []chan string) error
	inputName   string
	outputNames []string
}

type multiplexerDesc struct {
	fn         func(ctx context.Context, inputs []chan string, output chan string) error
	inputNames []string
	outputName string
}

type Conveyer interface {
	RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error,
		input string, output string)
	RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string, output string)
	RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error,
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
		mu:      sync.RWMutex{},
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
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string, output string,
) {
	c.decs = append(c.decs, decoratorDesc{
		fn:         fn,
		inputName:  input,
		outputName: output,
	})
}

func (c *conveyorImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string, outputs []string,
) {
	c.seps = append(c.seps, separatorDesc{
		fn:          fn,
		inputName:   input,
		outputNames: outputs,
	})
}

func (c *conveyorImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string, output string,
) {
	c.muxes = append(c.muxes, multiplexerDesc{
		fn:         fn,
		inputNames: inputs,
		outputName: output,
	})
}
func (c *conveyorImpl) runDecorators(ctx context.Context, waitGroup *sync.WaitGroup) {
	for _, decorator := range c.decs {
		inputCh := c.getOrCreate(decorator.inputName)
		outputCh := c.getOrCreate(decorator.outputName)

		waitGroup.Add(1)

		go func(desc decoratorDesc, input, output chan string) {
			defer waitGroup.Done()

			err := desc.fn(ctx, input, output)

			if err != nil {
				c.errCh <- err
			}
		}(decorator, inputCh, outputCh)
	}
}

func (c *conveyorImpl) runSeparators(ctx context.Context, waitGroup *sync.WaitGroup) {
	for _, separator := range c.seps {
		inputCh := c.getOrCreate(separator.inputName)
		outs := make([]chan string, 0, len(separator.outputNames))

		for _, name := range separator.outputNames {
			outs = append(outs, c.getOrCreate(name))
		}

		waitGroup.Add(1)

		go func(desc separatorDesc, input chan string, outputs []chan string) {
			defer waitGroup.Done()

			err := desc.fn(ctx, input, outputs)

			if err != nil {
				c.errCh <- err
			}
		}(separator, inputCh, outs)
	}
}

func (c *conveyorImpl) runMultiplexers(ctx context.Context, waitGroup *sync.WaitGroup) {
	for _, multiplexer := range c.muxes {
		inputs := make([]chan string, 0, len(multiplexer.inputNames))
		for _, name := range multiplexer.inputNames {
			inputs = append(inputs, c.getOrCreate(name))
		}

		outputCh := c.getOrCreate(multiplexer.outputName)

		waitGroup.Add(1)

		go func(desc multiplexerDesc, ins []chan string, output chan string) {
			defer waitGroup.Done()

			err := desc.fn(ctx, ins, output)

			if err != nil {
				c.errCh <- err
			}
		}(multiplexer, inputs, outputCh)
	}
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

	c.runDecorators(ctx, &waitGroup)
	c.runSeparators(ctx, &waitGroup)
	c.runMultiplexers(ctx, &waitGroup)

	go func() {
		waitGroup.Wait()
		close(c.errCh)
	}()

	return nil
}

func (c *conveyorImpl) Send(input string, data string) error {
	c.mu.RLock()
	channel, channelExists := c.chans[input]
	c.mu.RUnlock()

	if !channelExists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *conveyorImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, channelExists := c.chans[output]
	c.mu.RUnlock()

	if !channelExists {
		return "", ErrChanNotFound
	}

	value, isOpen := <-channel
	if !isOpen {
		return "", ErrClosed
	}

	return value, nil
}
