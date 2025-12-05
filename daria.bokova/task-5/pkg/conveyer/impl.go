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
	errChan      chan error
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

func newConveyerImpl(size int) *conveyerImpl {
	return &conveyerImpl{
		channels:     make(map[string]chan string),
		size:         size,
		decorators:   make([]decoratorHandler, 0),
		multiplexers: make([]multiplexerHandler, 0),
		separators:   make([]separatorHandler, 0),
		errChan:      make(chan error, 10),
	}
}

func (c *conveyerImpl) RegisterDecorator(fn DecoratorFunc, input string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Создаем каналы если их нет
	c.ensureChannel(input)
	c.ensureChannel(output)

	c.decorators = append(c.decorators, decoratorHandler{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *conveyerImpl) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Создаем каналы если их нет
	for _, input := range inputs {
		c.ensureChannel(input)
	}
	c.ensureChannel(output)

	c.multiplexers = append(c.multiplexers, multiplexerHandler{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *conveyerImpl) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Создаем каналы если их нет
	c.ensureChannel(input)
	for _, output := range outputs {
		c.ensureChannel(output)
	}

	c.separators = append(c.separators, separatorHandler{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *conveyerImpl) ensureChannel(name string) {
	if _, exists := c.channels[name]; !exists {
		c.channels[name] = make(chan string, c.size)
	}
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return errors.New("conveyer is already running")
	}
	c.running = true

	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.mu.Unlock()

	defer func() {
		cancel()
		c.stop()
	}()

	// Запускаем все обработчики
	c.startHandlers(runCtx)

	// Ждем завершения
	select {
	case <-runCtx.Done():
		return nil
	case err := <-c.errChan:
		return err
	}
}

func (c *conveyerImpl) startHandlers(ctx context.Context) {
	// Декораторы
	for _, h := range c.decorators {
		c.wg.Add(1)
		go func(h decoratorHandler) {
			defer c.wg.Done()

			inputChan := c.getChannel(h.input)
			outputChan := c.getChannel(h.output)

			if err := h.fn(ctx, inputChan, outputChan); err != nil {
				select {
				case c.errChan <- err:
				default:
				}
			}
		}(h)
	}

	// Мультиплексоры
	for _, h := range c.multiplexers {
		c.wg.Add(1)
		go func(h multiplexerHandler) {
			defer c.wg.Done()

			inputChans := make([]chan string, len(h.inputs))
			for i, inputName := range h.inputs {
				inputChans[i] = c.getChannel(inputName)
			}
			outputChan := c.getChannel(h.output)

			if err := h.fn(ctx, inputChans, outputChan); err != nil {
				select {
				case c.errChan <- err:
				default:
				}
			}
		}(h)
	}

	// Сепараторы
	for _, h := range c.separators {
		c.wg.Add(1)
		go func(h separatorHandler) {
			defer c.wg.Done()

			inputChan := c.getChannel(h.input)
			outputChans := make([]chan string, len(h.outputs))
			for i, outputName := range h.outputs {
				outputChans[i] = c.getChannel(outputName)
			}

			if err := h.fn(ctx, inputChan, outputChans); err != nil {
				select {
				case c.errChan <- err:
				default:
				}
			}
		}(h)
	}
}

func (c *conveyerImpl) stop() {
	c.mu.Lock()
	c.running = false
	channels := make(map[string]chan string)
	for k, v := range c.channels {
		channels[k] = v
	}
	c.mu.Unlock()

	// Закрываем каналы
	for _, ch := range channels {
		close(ch)
	}

	c.wg.Wait()
}

func (c *conveyerImpl) getChannel(name string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channels[name]
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch := c.getChannel(input)
	if ch == nil {
		return errors.New("chan not found")
	}

	// Неблокирующая отправка
	select {
	case ch <- data:
		return nil
	default:
		// Канал заполнен
		return nil
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch := c.getChannel(output)
	if ch == nil {
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
