package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrUndefined    = errors.New("undefined")
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
	size int

	mu      sync.RWMutex
	chans   map[string]chan string
	decs    []decoratorDesc
	seps    []separatorDesc
	muxes   []multiplexerDesc
	errCh   chan error
	started bool
}

func New(size int) Conveyer {
	return &conveyorImpl{
		size:  size,
		chans: make(map[string]chan string),
		errCh: make(chan error, size),
	}
}

func (c *conveyorImpl) getOrCreate(name string) chan string {
	ch, ok := c.chans[name]
	if !ok {
		ch = make(chan string, c.size)
		c.chans[name] = ch
	}
	return ch
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
		return fmt.Errorf("already started")
	}
	c.started = true
	c.mu.Unlock()

	var wg sync.WaitGroup

	for _, d := range c.decs {
		in := c.getOrCreate(d.inputName)
		out := c.getOrCreate(d.outputName)

		wg.Add(1)
		go func(desc decoratorDesc, in, out chan string) {
			defer wg.Done()
			desc.fn(ctx, in, out, c.errCh)
		}(d, in, out)
	}

	for _, s := range c.seps {
		in := c.getOrCreate(s.inputName)
		outs := make([]chan string, 0, len(s.outputNames))
		for _, name := range s.outputNames {
			outs = append(outs, c.getOrCreate(name))
		}

		wg.Add(1)
		go func(desc separatorDesc, in chan string, outs []chan string) {
			defer wg.Done()
			desc.fn(ctx, in, outs, c.errCh)
		}(s, in, outs)
	}

	for _, m := range c.muxes {
		inputs := make([]chan string, 0, len(m.inputNames))
		for _, name := range m.inputNames {
			inputs = append(inputs, c.getOrCreate(name))
		}
		out := c.getOrCreate(m.outputName)

		wg.Add(1)
		go func(desc multiplexerDesc, ins []chan string, out chan string) {
			defer wg.Done()
			desc.fn(ctx, ins, out, c.errCh)
		}(m, inputs, out)
	}

	go func() {
		wg.Wait()
		close(c.errCh)
	}()

	return nil
}

func (c *conveyorImpl) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.chans[input]
	c.mu.RUnlock()
	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *conveyorImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.chans[output]
	c.mu.RUnlock()
	if !ok {
		return "", ErrChanNotFound
	}

	v, ok := <-ch
	if !ok {
		return "", fmt.Errorf("closed")
	}
	return v, nil
}
