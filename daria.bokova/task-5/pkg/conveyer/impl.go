package conveyer

import (
	"context"
	"errors"
	"fmt"
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
	ctx          context.Context
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
	// Создаем каналы при регистрации
	c.getOrCreateChannel(input)
	c.getOrCreateChannel(output)

	c.decorators = append(c.decorators, decoratorHandler{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *conveyerImpl) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	// Создаем каналы при регистрации
	for _, input := range inputs {
		c.getOrCreateChannel(input)
	}
	c.getOrCreateChannel(output)

	c.multiplexers = append(c.multiplexers, multiplexerHandler{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *conveyerImpl) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	// Создаем каналы при регистрации
	c.getOrCreateChannel(input)
	for _, output := range outputs {
		c.getOrCreateChannel(output)
	}

	c.separators = append(c.separators, separatorHandler{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return errors.New("conveyer is already running")
	}

	c.running = true
	c.ctx, c.cancel = context.WithCancel(ctx)
	ctx = c.ctx
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
		c.closeAllChannels()
		c.wg.Wait()
	}()

	errChan := make(chan error, len(c.decorators)+len(c.multiplexers)+len(c.separators))

	// Запускаем декораторы
	for _, d := range c.decorators {
		inputChan := c.getChannel(d.input)
		outputChan := c.getChannel(d.output)

		c.wg.Add(1)
		go func(d decoratorHandler) {
			defer c.wg.Done()
			if err := d.fn(ctx, inputChan, outputChan); err != nil {
				select {
				case errChan <- err:
				default:
				}
				c.cancel()
			}
		}(d)
	}

	// Запускаем мультиплексоры
	for _, m := range c.multiplexers {
		inputChans := make([]chan string, len(m.inputs))
		for i, inputName := range m.inputs {
			inputChans[i] = c.getChannel(inputName)
		}
		outputChan := c.getChannel(m.output)

		c.wg.Add(1)
		go func(m multiplexerHandler) {
			defer c.wg.Done()
			if err := m.fn(ctx, inputChans, outputChan); err != nil {
				select {
				case errChan <- err:
				default:
				}
				c.cancel()
			}
		}(m)
	}

	// Запускаем сепараторы
	for _, s := range c.separators {
		inputChan := c.getChannel(s.input)
		outputChans := make([]chan string, len(s.outputs))
		for i, outputName := range s.outputs {
			outputChans[i] = c.getChannel(outputName)
		}

		c.wg.Add(1)
		go func(s separatorHandler) {
			defer c.wg.Done()
			if err := s.fn(ctx, inputChan, outputChans); err != nil {
				select {
				case errChan <- err:
				default:
				}
				c.cancel()
			}
		}(s)
	}

	// Ждем завершения контекста или ошибок
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func (c *conveyerImpl) getChannel(name string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channels[name]
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		select {
		case _, ok := <-ch:
			if !ok {
				// Канал уже закрыт
				continue
			}
		default:
		}
		close(ch)
		delete(c.channels, name)
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	ch, exists := c.channels[input]
	running := c.running
	c.mu.RUnlock()

	if !exists {
		return errors.New("chan not found")
	}

	if !running {
		return errors.New("conveyer is not running")
	}

	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case ch <- data:
		return nil
	default:
		return fmt.Errorf("channel %s is full", input)
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
