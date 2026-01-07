package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	errAlreadyRunning = errors.New("already running")
	errChanNotFound   = errors.New("chan not found")
)

type Conveyer struct {
	size       int
	channels   map[string]chan string
	handlers   []func(ctx context.Context) error
	mutex      sync.RWMutex
	isRunning  bool
	cancelFunc context.CancelFunc
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:      size,
		channels:  make(map[string]chan string),
		handlers:  []func(ctx context.Context) error{},
		mutex:     sync.RWMutex{},
		isRunning: false,
	}
}

func (c *Conveyer) getChannel(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunction func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inputChannel := c.getChannel(input)
	outputChannel := c.getChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return decoratorFunction(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunction func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))
	for index, name := range inputs {
		inputChannels[index] = c.getChannel(name)
	}

	outputChannel := c.getChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return multiplexerFunction(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunction func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inputChannel := c.getChannel(input)
	outputChannels := make([]chan string, len(outputs))

	for index, name := range outputs {
		outputChannels[index] = c.getChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return separatorFunction(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mutex.Lock()

	if c.isRunning {
		c.mutex.Unlock()
		return errAlreadyRunning
	}

	c.isRunning = true

	ctx, cancel := context.WithCancel(ctx)
	c.cancelFunc = cancel

	c.mutex.Unlock()

	errorGroup, ctxWithCancel := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		handlerCopy := handler
		errorGroup.Go(func() error {
			return handlerCopy(ctxWithCancel)
		})
	}

	resultChan := make(chan error, 1)
	go func() {
		resultChan <- errorGroup.Wait()
		close(resultChan)
	}()

	select {
	case <-ctx.Done():
		cancel()
		return fmt.Errorf("context canceled: %w", ctx.Err())
	case err := <-resultChan:
		cancel()
		return err
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mutex.RLock()
	channel, exists := c.channels[input]
	c.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("%w", errChanNotFound)
	}

	channel <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mutex.RLock()
	channel, exists := c.channels[output]
	c.mutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("%w", errChanNotFound)
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
