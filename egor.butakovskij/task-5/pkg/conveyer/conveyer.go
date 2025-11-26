package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	errChanNotFound   = errors.New("chan not found")
	ErrInvalidHandler = errors.New("invalid handler function signature")
)

type Conveyer struct {
	channels       map[string]chan string
	HandlerConfigs []HandlerConfig
	bufferSize     int
	mu             sync.Mutex
	wg             sync.WaitGroup
}

type HandlerConfig struct {
	Fn        interface{}
	InputIds  []string
	OutputIds []string
	Type      int
}

const (
	DecoratorType   = 1
	MultiplexerType = 2
	SeparatorType   = 3
)

func New(size int) *Conveyer {
	conv := &Conveyer{
		channels:       make(map[string]chan string),
		HandlerConfigs: make([]HandlerConfig, 0),
		bufferSize:     size,
		mu:             sync.Mutex{},
		wg:             sync.WaitGroup{},
	}
	return conv
}

func (c *Conveyer) getChan(name string) chan string {
	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getChan(input)
	c.getChan(output)

	config := HandlerConfig{
		Type:      DecoratorType,
		Fn:        fn,
		InputIds:  []string{input},
		OutputIds: []string{output},
	}

	c.HandlerConfigs = append(c.HandlerConfigs, config)
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, input := range inputs {
		c.getChan(input)
	}

	c.getChan(output)

	config := HandlerConfig{
		Type:      MultiplexerType,
		Fn:        fn,
		InputIds:  inputs,
		OutputIds: []string{output},
	}

	c.HandlerConfigs = append(c.HandlerConfigs, config)
}

func (c *Conveyer) RegisterSeparator(
	fn func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, output := range outputs {
		c.getChan(output)
	}

	c.getChan(input)

	config := HandlerConfig{
		Type:      SeparatorType,
		Fn:        fn,
		InputIds:  []string{input},
		OutputIds: outputs,
	}

	c.HandlerConfigs = append(c.HandlerConfigs, config)
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	ch, exists := c.channels[input]
	c.mu.Unlock()
	if !exists {
		return errChanNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	ch, exists := c.channels[output]
	c.mu.Unlock()
	if !exists {
		return "", errChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (c *Conveyer) runHandler(cfg HandlerConfig, inputs []chan string, outputs []chan string, errChan chan error, ctx context.Context) {
	defer c.wg.Done()
	var err error
	switch cfg.Type {
	case DecoratorType:

		fn, ok := cfg.Fn.(func(context.Context, chan string, chan string) error)
		if !ok {
			errChan <- ErrInvalidHandler
			return
		}

		err = fn(ctx, inputs[0], outputs[0])
	case MultiplexerType:
		fn, ok := cfg.Fn.(func(context.Context, []chan string, chan string) error)
		if !ok {
			errChan <- ErrInvalidHandler
			return
		}

		err = fn(ctx, inputs, outputs[0])
	case SeparatorType:

		fn, ok := cfg.Fn.(func(context.Context, chan string, []chan string) error)
		if !ok {
			errChan <- ErrInvalidHandler
			return
		}

		err = fn(ctx, inputs[0], outputs)
	}

	if err != nil {
		errChan <- err
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	errChan := make(chan error, len(c.HandlerConfigs))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c.mu.Lock()

	for _, config := range c.HandlerConfigs {
		var (
			inputChans  []chan string
			outputChans []chan string
		)

		for _, id := range config.InputIds {
			inputChans = append(inputChans, c.channels[id])
		}

		for _, id := range config.OutputIds {
			outputChans = append(outputChans, c.channels[id])
		}

		c.wg.Add(1)

		go c.runHandler(config, inputChans, outputChans, errChan, ctx)
	}
	c.mu.Unlock()
	select {
	case err := <-errChan:
		cancel()
		c.wg.Wait()
		return err
	case <-ctx.Done():
		c.wg.Wait()
		return nil
	}
}
