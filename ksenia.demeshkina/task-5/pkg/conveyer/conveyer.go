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
	cancel         context.CancelFunc
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
			return "", errors.New("channel closed")
		}
		return data, nil
	default:
		return "", errors.New("no data available")
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	for _, config := range c.HandlerConfigs {
		c.wg.Add(1)
		go func(cfg HandlerConfig) {
			defer c.wg.Done()

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

			var err error
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

			c.mu.Lock()
			for _, outID := range cfg.OutputIds {
				ch, exists := c.channels[outID]
				if exists && ch != nil {
					close(ch)
					delete(c.channels, outID)
				}
			}
			c.mu.Unlock()

			if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
				select {
				case errChan <- err:
					cancel()
				default:
				}
			}
		}(config)
	}

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:
		c.stopAll()
		return err
	case <-done:
		c.stopAll()
		return nil
	case <-ctx.Done():
		c.stopAll()
		return ctx.Err()
	}
}

func (c *Conveyer) stopAll() {
	c.stopOnce.Do(func() {
		if c.cancel != nil {
			c.cancel()
		}

		c.mu.Lock()
		for name, ch := range c.channels {
			if ch != nil {
				close(ch)
			}
			delete(c.channels, name)
		}
		c.mu.Unlock()

		c.wg.Wait()
	})
}
