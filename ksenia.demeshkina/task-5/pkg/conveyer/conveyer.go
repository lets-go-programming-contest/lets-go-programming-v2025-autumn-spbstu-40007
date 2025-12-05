package conveyer

import (
	"context"
	"errors"
	"sync"
)

var errChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channels       map[string]chan string
	HandlerConfigs []HandlerConfig
	bufferSize     int
	mu             sync.Mutex
	wg             sync.WaitGroup
	stopOnce       sync.Once
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
	return &Conveyer{
		channels:       make(map[string]chan string),
		HandlerConfigs: make([]HandlerConfig, 0),
		bufferSize:     size,
		mu:             sync.Mutex{},
		wg:             sync.WaitGroup{},
		stopOnce:       sync.Once{},
	}
}

func (c *Conveyer) getOrCreateChan(name string) chan string {
	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fun func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChan(input)
	c.getOrCreateChan(output)

	config := HandlerConfig{
		Type:      DecoratorType,
		Fn:        fun,
		InputIds:  []string{input},
		OutputIds: []string{output},
	}

	c.HandlerConfigs = append(c.HandlerConfigs, config)
}

func (c *Conveyer) RegisterMultiplexer(
	fun func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, input := range inputs {
		c.getOrCreateChan(input)
	}
	c.getOrCreateChan(output)

	config := HandlerConfig{
		Type:      MultiplexerType,
		Fn:        fun,
		InputIds:  inputs,
		OutputIds: []string{output},
	}

	c.HandlerConfigs = append(c.HandlerConfigs, config)
}

func (c *Conveyer) RegisterSeparator(
	fun func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChan(input)
	for _, output := range outputs {
		c.getOrCreateChan(output)
	}

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

	select {
	case chann <- data:
		return nil
	default:
		return errors.New("channel is full")
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	chann, exists := c.channels[output]
	c.mu.Unlock()

	if !exists {
		return "", errChanNotFound
	}

	select {
	case data, ok := <-chann:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	default:
		return "", errors.New("no data available")
	}
}

func (c *Conveyer) Hndl(cfg HandlerConfig, ctx context.Context, errC chan error) {
	defer c.wg.Done()

	var err error
	c.mu.Lock()

	inputChans := make([]chan string, len(cfg.InputIds))
	for i, id := range cfg.InputIds {
		inputChans[i] = c.channels[id]
	}

	outputChans := make([]chan string, len(cfg.OutputIds))
	for i, id := range cfg.OutputIds {
		outputChans[i] = c.channels[id]
	}
	c.mu.Unlock()

	switch cfg.Type {
	case DecoratorType:
		fun, ok := cfg.Fn.(func(context.Context, chan string, chan string) error)
		if !ok {
			return
		}
		err = fun(ctx, inputChans[0], outputChans[0])
	case MultiplexerType:
		fun, ok := cfg.Fn.(func(context.Context, []chan string, chan string) error)
		if !ok {
			return
		}
		err = fun(ctx, inputChans, outputChans[0])
	case SeparatorType:
		fun, ok := cfg.Fn.(func(context.Context, chan string, []chan string) error)
		if !ok {
			return
		}
		err = fun(ctx, inputChans[0], outputChans)
	}

	if err != nil {
		select {
		case errC <- err:
		default:
		}
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	errChan := make(chan error, len(c.HandlerConfigs))

	for _, config := range c.HandlerConfigs {
		c.wg.Add(1)
		go c.Hndl(config, ctx, errChan)
	}

	go func() {
		c.wg.Wait()
		close(errChan)
	}()

	select {
	case err := <-errChan:
		c.stopAll()
		return err
	case <-ctx.Done():
		c.stopAll()
		return ctx.Err()
	}
}

func (c *Conveyer) stopAll() {
	c.stopOnce.Do(func() {
		c.mu.Lock()
		closedChannels := make(map[chan string]bool)
		for name, ch := range c.channels {
			if !closedChannels[ch] {
				closedChannels[ch] = true
				close(ch)
			}
			delete(c.channels, name)
		}
		c.mu.Unlock()
		c.wg.Wait()
	})
}
