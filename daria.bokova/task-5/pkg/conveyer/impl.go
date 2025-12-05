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

func (c *conveyerImpl) getChannel(name string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channels[name]
}

func (c *conveyerImpl) RegisterDecorator(fn DecoratorFunc, input string, output string) {
	c.getOrCreateChannel(input)
	c.getOrCreateChannel(output)

	c.decorators = append(c.decorators, decoratorHandler{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *conveyerImpl) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
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

	// Создаем новый контекст с отменой
	var cancelCtx context.Context
	cancelCtx, c.cancel = context.WithCancel(ctx)
	c.mu.Unlock()

	// Канал для ошибок
	errChan := make(chan error, 1)

	// Запускаем обработчики
	for _, d := range c.decorators {
		c.wg.Add(1)
		go func(d decoratorHandler) {
			defer c.wg.Done()
			inputChan := c.getChannel(d.input)
			outputChan := c.getChannel(d.output)

			if err := d.fn(cancelCtx, inputChan, outputChan); err != nil {
				select {
				case errChan <- err:
				default:
				}
				c.cancel()
			}
		}(d)
	}

	for _, m := range c.multiplexers {
		c.wg.Add(1)
		go func(m multiplexerHandler) {
			defer c.wg.Done()
			inputChans := make([]chan string, len(m.inputs))
			for i, inputName := range m.inputs {
				inputChans[i] = c.getChannel(inputName)
			}
			outputChan := c.getChannel(m.output)

			if err := m.fn(cancelCtx, inputChans, outputChan); err != nil {
				select {
				case errChan <- err:
				default:
				}
				c.cancel()
			}
		}(m)
	}

	for _, s := range c.separators {
		c.wg.Add(1)
		go func(s separatorHandler) {
			defer c.wg.Done()
			inputChan := c.getChannel(s.input)
			outputChans := make([]chan string, len(s.outputs))
			for i, outputName := range s.outputs {
				outputChans[i] = c.getChannel(outputName)
			}

			if err := s.fn(cancelCtx, inputChan, outputChans); err != nil {
				select {
				case errChan <- err:
				default:
				}
				c.cancel()
			}
		}(s)
	}

	// Ждем завершения
	select {
	case <-cancelCtx.Done():
		c.stop()
		return cancelCtx.Err()
	case err := <-errChan:
		c.stop()
		return err
	}
}

func (c *conveyerImpl) stop() {
	c.mu.Lock()
	if c.cancel != nil {
		c.cancel()
	}
	c.running = false
	c.mu.Unlock()

	// Закрываем все каналы
	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.channels = make(map[string]chan string)
	c.mu.Unlock()

	// Ждем завершения всех горутин
	c.wg.Wait()
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch := c.getChannel(input)
	if ch == nil {
		return errors.New("chan not found")
	}

	select {
	case ch <- data:
		return nil
	default:
		// Канал заполнен - игнорируем
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
