package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound    = errors.New("chan not found")
	ErrConveyerStopped = errors.New("conveyer stopped")
	ErrConveyerRunning = errors.New("conveyer already running")
	ErrNotRunning      = errors.New("conveyer not running")
	ErrBufferFull      = errors.New("channel buffer full")
	ErrNoData          = errors.New("no data available")
)

type conveyerImpl struct {
	mu          sync.RWMutex
	channels    map[string]chan string
	bufferSize  int
	handlers    []handlerInfo
	cancelFuncs []context.CancelFunc
	running     bool
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

type handlerInfo struct {
	handlerType string
	fn          interface{}
	inputs      []string
	outputs     []string
}

func New(size int) Conveyer {
	return &conveyerImpl{
		channels:   make(map[string]chan string),
		bufferSize: size,
		stopChan:   make(chan struct{}),
	}
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) getChannel(name string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if ch, exists := c.channels[name]; exists {
		return ch, nil
	}
	return nil, ErrChanNotFound
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, handlerInfo{
		handlerType: "decorator",
		fn:          fn,
		inputs:      []string{input},
		outputs:     []string{output},
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, handlerInfo{
		handlerType: "multiplexer",
		fn:          fn,
		inputs:      inputs,
		outputs:     []string{output},
	})
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, handlerInfo{
		handlerType: "separator",
		fn:          fn,
		inputs:      []string{input},
		outputs:     outputs,
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return ErrConveyerRunning
	}
	c.running = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		c.closeAllChannels()
		c.mu.Unlock()
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, len(c.handlers))

	for _, info := range c.handlers {
		c.wg.Add(1)
		info := info
		go func() {
			defer c.wg.Done()

			handlerCtx, handlerCancel := context.WithCancel(ctx)
			c.mu.Lock()
			c.cancelFuncs = append(c.cancelFuncs, handlerCancel)
			c.mu.Unlock()

			var err error
			switch info.handlerType {
			case "decorator":
				fn := info.fn.(func(ctx context.Context, input chan string, output chan string) error)
				inputCh := c.getOrCreateChannel(info.inputs[0])
				outputCh := c.getOrCreateChannel(info.outputs[0])
				err = fn(handlerCtx, inputCh, outputCh)
			case "multiplexer":
				fn := info.fn.(func(ctx context.Context, inputs []chan string, output chan string) error)
				inputChs := make([]chan string, len(info.inputs))
				for i, inputName := range info.inputs {
					inputChs[i] = c.getOrCreateChannel(inputName)
				}
				outputCh := c.getOrCreateChannel(info.outputs[0])
				err = fn(handlerCtx, inputChs, outputCh)
			case "separator":
				fn := info.fn.(func(ctx context.Context, input chan string, outputs []chan string) error)
				inputCh := c.getOrCreateChannel(info.inputs[0])
				outputChs := make([]chan string, len(info.outputs))
				for i, outputName := range info.outputs {
					outputChs[i] = c.getOrCreateChannel(outputName)
				}
				err = fn(handlerCtx, inputCh, outputChs)
			}

			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				cancel()
			}
		}()
	}

	go func() {
		c.wg.Wait()
		close(errChan)
	}()

	select {
	case err, ok := <-errChan:
		if ok && err != nil {
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-c.stopChan:
		return ErrConveyerStopped
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[input]
	if !exists {
		return ErrChanNotFound
	}

	if !c.running {
		return ErrNotRunning
	}

	select {
	case ch <- data:
		return nil
	default:
		return ErrBufferFull
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, err := c.getChannel(output)
	if err != nil {
		return "", err
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	default:
		return "", ErrNoData
	}
}

func (c *conveyerImpl) Stop() {
	close(c.stopChan)
}

func (c *conveyerImpl) closeAllChannels() {
	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}

	for _, cancel := range c.cancelFuncs {
		cancel()
	}
	c.cancelFuncs = nil
}
