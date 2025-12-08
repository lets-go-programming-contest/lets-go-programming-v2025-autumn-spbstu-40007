package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound       = errors.New("chan not found")
	ErrInvalidDecorator   = errors.New("invalid decorator function type")
	ErrInvalidMultiplexer = errors.New("invalid multiplexer function type")
	ErrInvalidSeparator   = errors.New("invalid separator function type")
)

const (
	typeDecorator   = 1
	typeMultiplexer = 2
	typeSeparator   = 3
)

type handlerConfig struct {
	hType   int
	inputs  []string
	outputs []string
	fn      interface{}
}

type Conveyer struct {
	mu         sync.Mutex
	channels   map[string]chan string
	configs    []handlerConfig
	bufferSize int
	waitGroup  sync.WaitGroup
	done       chan struct{}
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:         sync.Mutex{},
		channels:   make(map[string]chan string),
		configs:    []handlerConfig{},
		bufferSize: size,
		waitGroup:  sync.WaitGroup{},
		done:       make(chan struct{}),
	}
}

func (c *Conveyer) getOrCreateChan(name string) {
	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, c.bufferSize)
	}
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChan(input)
	c.getOrCreateChan(output)

	c.configs = append(c.configs, handlerConfig{
		hType:   typeDecorator,
		inputs:  []string{input},
		outputs: []string{output},
		fn:      handlerFunc,
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, inputChan := range inputs {
		c.getOrCreateChan(inputChan)
	}

	c.getOrCreateChan(output)

	c.configs = append(c.configs, handlerConfig{
		hType:   typeMultiplexer,
		inputs:  inputs,
		outputs: []string{output},
		fn:      handlerFunc,
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChan(input)

	for _, outputChan := range outputs {
		c.getOrCreateChan(outputChan)
	}

	c.configs = append(c.configs, handlerConfig{
		hType:   typeSeparator,
		inputs:  []string{input},
		outputs: outputs,
		fn:      handlerFunc,
	})
}

func (c *Conveyer) runHandler(ctx context.Context, conf handlerConfig, errChan chan error, cancel context.CancelFunc) {
	defer c.waitGroup.Done()

	c.mu.Lock()

	inputChannels := make([]chan string, 0, len(conf.inputs))
	for _, name := range conf.inputs {
		inputChannels = append(inputChannels, c.channels[name])
	}

	outputChannels := make([]chan string, 0, len(conf.outputs))
	for _, name := range conf.outputs {
		outputChannels = append(outputChannels, c.channels[name])
	}

	c.mu.Unlock()

	var err error

	switch conf.hType {
	case typeDecorator:
		fn, ok := conf.fn.(func(context.Context, chan string, chan string) error)
		if !ok {
			err = ErrInvalidDecorator
		} else {
			err = fn(ctx, inputChannels[0], outputChannels[0])
		}

	case typeMultiplexer:
		fn, ok := conf.fn.(func(context.Context, []chan string, chan string) error)
		if !ok {
			err = ErrInvalidMultiplexer
		} else {
			err = fn(ctx, inputChannels, outputChannels[0])
		}

	case typeSeparator:
		fn, ok := conf.fn.(func(context.Context, chan string, []chan string) error)
		if !ok {
			err = ErrInvalidSeparator
		} else {
			err = fn(ctx, inputChannels[0], outputChannels)
		}
	}

	if err != nil {
		select {
		case errChan <- err:
			cancel()
		default:
		}
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 1)

	c.done = make(chan struct{})

	for _, cfg := range c.configs {
		c.waitGroup.Add(1)

		go c.runHandler(ctx, cfg, errChan, cancel)
	}

	go func() {
		c.waitGroup.Wait()
		close(c.done)
	}()

	var resultErr error

	select {
	case err := <-errChan:
		resultErr = err

	case <-ctx.Done():
		resultErr = nil

	case <-c.done:
		resultErr = nil
	}

	<-c.done

	c.mu.Lock()
	for _, channel := range c.channels {
		close(channel)
	}
	c.mu.Unlock()

	return resultErr
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	channel, ok := c.channels[input]
	c.mu.Unlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	channel, ok := c.channels[output]
	c.mu.Unlock()

	if !ok {
		return "", ErrChanNotFound
	}

	data, open := <-channel
	if !open {
		return "undefined", nil
	}

	return data, nil
}
