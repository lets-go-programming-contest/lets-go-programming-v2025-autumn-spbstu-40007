package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer interface {
	RegisterDecorator(
		handlerFuncParam func(context.Context, chan string, chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		handlerFuncParam func(context.Context, []chan string, chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		handlerFuncParam func(context.Context, chan string, []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type handlerRunner func(context.Context) error

type conveyer struct {
	size     int
	mu       sync.RWMutex
	chans    map[string]chan string
	handlers []handlerRunner
}

func New(size int) Conveyer {
	return &conveyer{
		size:     size,
		mu:       sync.RWMutex{},
		chans:    make(map[string]chan string),
		handlers: nil,
	}
}

func (c *conveyer) getOrCreate(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, ok := c.chans[name]; ok {

		return channel
	}

	channel := make(chan string, c.size)
	c.chans[name] = channel

	return channel
}

func (c *conveyer) RegisterDecorator(
	handlerFuncParam func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inputChan := c.getOrCreate(input)
	outputChan := c.getOrCreate(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {

		return handlerFuncParam(ctx, inputChan, outputChan)
	})
}

func (c *conveyer) RegisterMultiplexer(
	handlerFuncParam func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, 0, len(inputs))
	for _, inputName := range inputs {
		inputChans = append(inputChans, c.getOrCreate(inputName))
	}

	outputChan := c.getOrCreate(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {

		return handlerFuncParam(ctx, inputChans, outputChan)
	})
}

func (c *conveyer) RegisterSeparator(
	handlerFuncParam func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := c.getOrCreate(input)

	outputChans := make([]chan string, 0, len(outputs))
	for _, outputName := range outputs {
		outputChans = append(outputChans, c.getOrCreate(outputName))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {

		return handlerFuncParam(ctx, inputChan, outputChans)
	})
}

func (c *conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	waitGroup := sync.WaitGroup{}
	errCh := make(chan error, 1)

	for _, handlerFunc := range c.handlers {
		handler := handlerFunc
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			if err := handler(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				cancel()
			}
		}()
	}

	select {
	case <-ctx.Done():

	case err := <-errCh:
		waitGroup.Wait()

		c.closeAll()

		return err
	}

	waitGroup.Wait()
	c.closeAll()

	return nil
}

func (c *conveyer) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range c.chans {
		close(channel)
	}
}

func (c *conveyer) Send(input string, data string) error {
	c.mu.RLock()
	inputChan, ok := c.chans[input]
	c.mu.RUnlock()

	if !ok {

		return ErrChanNotFound
	}

	inputChan <- data

	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	outputChan, ok := c.chans[output]
	c.mu.RUnlock()

	if !ok {

		return "", ErrChanNotFound
	}

	val, ok := <-outputChan
	if !ok {

		return "undefined", nil
	}

	return val, nil
}