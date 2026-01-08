package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const UndefinedResult = "undefined"

type Conveyer struct {
	size               int
	mutex              sync.RWMutex
	channels           map[string]chan string
	registeredHandlers []func(context.Context) error
}

func New(bufferSize int) *Conveyer {
	return &Conveyer{
		size:               bufferSize,
		mutex:              sync.RWMutex{},
		channels:           make(map[string]chan string),
		registeredHandlers: make([]func(context.Context) error, 0),
	}
}

func (c *Conveyer) getOrCreateChannel(channelName string) chan string {
	channel, channelExists := c.channels[channelName]
	if !channelExists {
		channel = make(chan string, c.size)
		c.channels[channelName] = channel
	}

	return channel
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(context.Context, chan string, chan string) error,
	inputChannelName string,
	outputChannelName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.getOrCreateChannel(inputChannelName)
	outputChannel := c.getOrCreateChannel(outputChannelName)

	c.registeredHandlers = append(c.registeredHandlers, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(context.Context, []chan string, chan string) error,
	inputChannelNames []string,
	outputChannelName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannels := make([]chan string, 0, len(inputChannelNames))
	for _, channelName := range inputChannelNames {
		inputChannels = append(inputChannels, c.getOrCreateChannel(channelName))
	}

	outputChannel := c.getOrCreateChannel(outputChannelName)

	c.registeredHandlers = append(c.registeredHandlers, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(context.Context, chan string, []chan string) error,
	inputChannelName string,
	outputChannelNames []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.getOrCreateChannel(inputChannelName)

	outputChannels := make([]chan string, 0, len(outputChannelNames))
	for _, channelName := range outputChannelNames {
		outputChannels = append(outputChannels, c.getOrCreateChannel(channelName))
	}

	c.registeredHandlers = append(c.registeredHandlers, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
}

func closeChannelSafely(channel chan string) {
	defer func() {
		_ = recover()
	}()
	close(channel)
}

func (c *Conveyer) closeAllChannelsSafely() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, channel := range c.channels {
		closeChannelSafely(channel)
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mutex.RLock()
	handlersCopy := make([]func(context.Context) error, len(c.registeredHandlers))
	copy(handlersCopy, c.registeredHandlers)
	c.mutex.RUnlock()

	group, contextWithCancel := errgroup.WithContext(ctx)

	for _, handler := range handlersCopy {
		specificHandler := handler

		group.Go(func() error {
			return specificHandler(contextWithCancel)
		})
	}

	err := group.Wait()

	c.closeAllChannelsSafely()

	if err != nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("conveyer failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(inputChannelName string, data string) error {
	channel, exists := func() (chan string, bool) {
		c.mutex.RLock()
		defer c.mutex.RUnlock()

		ch, ok := c.channels[inputChannelName]

		return ch, ok
	}()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(outputChannelName string) (string, error) {
	channel, exists := func() (chan string, bool) {
		c.mutex.RLock()
		defer c.mutex.RUnlock()

		ch, ok := c.channels[outputChannelName]

		return ch, ok
	}()

	if !exists {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return UndefinedResult, nil
	}

	return value, nil
}
