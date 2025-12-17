package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	errChanNotFound    = errors.New("chan not found")
	errDecoratorSign   = errors.New("invalid decorator signature")
	errMultiplexorSign = errors.New("invalid multiplexer signature")
	errSeparatorSign   = errors.New("invalid separator signature")
)

type handlerKind int

const (
	kindDecorator handlerKind = iota + 1
	kindMultiplexer
	kindSeparator
)

type handlerCfg struct {
	kind      handlerKind
	function  interface{}
	inputIDs  []string
	outputIDs []string
}

type Conveyer struct {
	channels map[string]chan string

	handlers []handlerCfg

	bufSize int
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		handlers: make([]handlerCfg, 0),
		bufSize:  size,
		mu:       sync.Mutex{},
		wg:       sync.WaitGroup{},
	}
}

func (c *Conveyer) ensureChan(id string) {
	_, ok := c.channels[id]
	if !ok {
		ch := make(chan string, c.bufSize)
		c.channels[id] = ch
	}
}

func (c *Conveyer) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureChan(input)
	c.ensureChan(output)

	cfg := handlerCfg{
		kind:      kindDecorator,
		function:  function,
		inputIDs:  []string{input},
		outputIDs: []string{output},
	}
	c.handlers = append(c.handlers, cfg)
}

func (c *Conveyer) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, id := range inputs {
		c.ensureChan(id)
	}

	c.ensureChan(output)

	cfg := handlerCfg{
		kind:      kindMultiplexer,
		function:  function,
		inputIDs:  append([]string(nil), inputs...),
		outputIDs: []string{output},
	}
	c.handlers = append(c.handlers, cfg)
}

func (c *Conveyer) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureChan(input)

	for _, id := range outputs {
		c.ensureChan(id)
	}

	cfg := handlerCfg{
		kind:      kindSeparator,
		function:  function,
		inputIDs:  []string{input},
		outputIDs: append([]string(nil), outputs...),
	}
	c.handlers = append(c.handlers, cfg)
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	channel, ok := c.channels[input]
	c.mu.Unlock()

	if !ok {
		return errChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	channel, ok := c.channels[output] //nolint:varnamelen
	c.mu.Unlock()

	if !ok {
		return "", errChanNotFound
	}

	v, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return v, nil
}

func (c *Conveyer) runHandler(ctx context.Context, cfg handlerCfg, errCh chan<- error) {
	defer c.wg.Done()

	c.mu.Lock()
	ins := make([]chan string, 0)
	outs := make([]chan string, 0)

	for _, id := range cfg.inputIDs {
		ins = append(ins, c.channels[id])
	}

	for _, id := range cfg.outputIDs {
		outs = append(outs, c.channels[id])
	}

	c.mu.Unlock()

	var err error

	switch cfg.kind {
	case kindDecorator:
		function, ok := cfg.function.(func(context.Context, chan string, chan string) error)
		if !ok {
			errCh <- errDecoratorSign

			return
		}

		err = function(ctx, ins[0], outs[0])

	case kindMultiplexer:
		function, ok := cfg.function.(func(context.Context, []chan string, chan string) error)
		if !ok {
			errCh <- errMultiplexorSign

			return
		}

		err = function(ctx, ins, outs[0])

	case kindSeparator:
		function, ok := cfg.function.(func(context.Context, chan string, []chan string) error)
		if !ok {
			errCh <- errSeparatorSign

			return
		}

		err = function(ctx, ins[0], outs)
	}

	if err != nil {
		select {
		case errCh <- err:
		default:
		}
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	if len(c.handlers) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	for _, cfg := range c.handlers {
		c.wg.Add(1)
		go c.runHandler(ctx, cfg, errCh)
	}

	select {
	case err := <-errCh:
		cancel()
		c.wg.Wait()

		return err
	case <-ctx.Done():
		c.wg.Wait()

		return nil
	}
}
