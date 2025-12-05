package conveyer

import (
	"context"
	"errors"
	"sync"
)

type conveyerImpl struct {
	mu           sync.RWMutex
	channels     map[string]chan string
	decorators   []decoratorHandler
	multiplexers []multiplexerHandler
	separators   []separatorHandler
	size         int
	running      bool
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

type decoratorHandler struct {
	fn     DecoratorFunc
	input  string
	output string
}

type multiplexerHandler struct {
	fn     MultiplexerFunc
	inputs []string
	output string
}

type separatorHandler struct {
	fn      SeparatorFunc
	input   string
	outputs []string
}

func New(size int) Conveyer {
	return &conveyerImpl{
		channels:     make(map[string]chan string),
		size:         size,
		decorators:   make([]decoratorHandler, 0),
		multiplexers: make([]multiplexerHandler, 0),
		separators:   make([]separatorHandler, 0),
	}
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) RegisterDecorator(fn DecoratorFunc, input string, output string) {
	c.decorators = append(c.decorators, decoratorHandler{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *conveyerImpl) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	c.multiplexers = append(c.multiplexers, multiplexerHandler{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *conveyerImpl) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	c.separators = append(c.separators, separatorHandler{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	if c.running {
		return errors.New("conveyer is already running")
	}

	c.mu.Lock()
	c.running = true
	ctx, c.cancel = context.WithCancel(ctx)
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
		c.closeAllChannels()
		c.wg.Wait()
	}()

	// Запускаем декораторы
	for _, d := range c.decorators {
		inputChan := c.getOrCreateChannel(d.input)
		outputChan := c.getOrCreateChannel(d.output)

		c.wg.Add(1)
		go func(d decoratorHandler) {
			defer c.wg.Done()
			if err := d.fn(ctx, inputChan, outputChan); err != nil {
				c.cancel()
			}
		}(d)
	}

	// Запускаем мультиплексоры
	for _, m := range c.multiplexers {
		inputChans := make([]chan string, len(m.inputs))
		for i, inputName := range m.inputs {
			inputChans[i] = c.getOrCreateChannel(inputName)
		}
		outputChan := c.getOrCreateChannel(m.output)

		c.wg.Add(1)
		go func(m multiplexerHandler) {
			defer c.wg.Done()
			if err := m.fn(ctx, inputChans, outputChan); err != nil {
				c.cancel()
			}
		}(m)
	}

	// Запускаем сепараторы
	for _, s := range c.separators {
		inputChan := c.getOrCreateChannel(s.input)
		outputChans := make([]chan string, len(s.outputs))
		for i, outputName := range s.outputs {
			outputChans[i] = c.getOrCreateChannel(outputName)
		}

		c.wg.Add(1)
		go func(s separatorHandler) {
			defer c.wg.Done()
			if err := s.fn(ctx, inputChan, outputChans); err != nil {
				c.cancel()
			}
		}(s)
	}

	// Ждем завершения контекста или ошибок
	<-ctx.Done()

	return ctx.Err()
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		select {
		case <-ch:
			// Канал уже закрыт или пуст
		default:
			close(ch)
		}
		delete(c.channels, name)
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	ch, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return errors.New("chan not found")
	}

	if !c.running {
		return errors.New("conveyer is not running")
	}

	select {
	case ch <- data:
		return nil
	default:
		return errors.New("channel is full")
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", errors.New("chan not found")
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	default:
		return "", errors.New("no data available")
	}
}
