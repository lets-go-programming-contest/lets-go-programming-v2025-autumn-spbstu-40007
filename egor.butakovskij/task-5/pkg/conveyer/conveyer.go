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

func getChan(name string, conv *Conveyer) {
	if _, exists := conv.channels[name]; exists {
		return
	}

	ch := make(chan string, conv.bufferSize)
	conv.channels[name] = ch
}

func (c *Conveyer) RegisterDecorator(
	fun func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	getChan(input, c)
	getChan(output, c)

	config := HandlerConfig{
		Type:      DecoratorType,
		Fn:        fun,
		InputIds:  []string{input},
		OutputIds: []string{output},
	}

	c.HandlerConfigs = append(c.HandlerConfigs, config)
}

func (c *Conveyer) RegisterMultiplexer(
	fun func(
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
		getChan(input, c)
	}

	getChan(output, c)

	config := HandlerConfig{
		Type:      MultiplexerType,
		Fn:        fun,
		InputIds:  inputs,
		OutputIds: []string{output},
	}

	c.HandlerConfigs = append(c.HandlerConfigs, config)
}

func (c *Conveyer) RegisterSeparator(
	fun func(
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
		getChan(output, c)
	}

	getChan(input, c)

	config := HandlerConfig{
		Type:      SeparatorType,
		Fn:        fun,
		InputIds:  []string{input},
		OutputIds: outputs,
	}

	c.HandlerConfigs = append(c.HandlerConfigs, config)
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	chann, exists := c.channels[input]
	c.mu.Unlock()

	if !exists {
		return errChanNotFound
	}

	chann <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	chann, exists := c.channels[output]
	c.mu.Unlock()

	if !exists {
		return "", errChanNotFound
	}

	data, ok := <-chann
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (c *Conveyer) Hndl(cfg HandlerConfig, ins []chan string, out []chan string, errC chan error, ctx context.Context) {
	defer c.wg.Done()

	var err error

	switch cfg.Type {
	case DecoratorType:
		fun, ok := cfg.Fn.(func(context.Context, chan string, chan string) error)
		if !ok {
			errC <- ErrInvalidHandler

			return
		}

		err = fun(ctx, ins[0], out[0])
	case MultiplexerType:
		fun, ok := cfg.Fn.(func(context.Context, []chan string, chan string) error)
		if !ok {
			errC <- ErrInvalidHandler

			return
		}

		err = fun(ctx, ins, out[0])
	case SeparatorType:
		fun, ok := cfg.Fn.(func(context.Context, chan string, []chan string) error)
		if !ok {
			errC <- ErrInvalidHandler

			return
		}

		err = fun(ctx, ins[0], out)
	}

	if err != nil {
		errC <- err
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

		go c.Hndl(config, inputChans, outputChans, errChan, ctx)
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
