package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChannelMissing = errors.New("chan not found")
	ErrAlreadyRunning = errors.New("conveyor already running")
)

const closedChannelValue = "undefined"

type handlerFn func(context.Context) error

type Conveyor struct {
	bufSize  int
	chans    map[string]chan string
	handlers []handlerFn
	mu       sync.Mutex
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(bufferSize int) *Conveyor {
	return &Conveyor{
		bufSize:  bufferSize,
		chans:    make(map[string]chan string),
		handlers: make([]handlerFn, 0),
		mu:       sync.Mutex{},
		wg:       sync.WaitGroup{},
		ctx:      nil,
		cancel:   nil,
	}
}

func (c *Conveyor) registerChannel(channelID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if existingChannel, exists := c.chans[channelID]; exists {
		return existingChannel
	}

	newChannel := make(chan string, c.bufSize)
	c.chans[channelID] = newChannel
	return newChannel
}

func (c *Conveyor) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	inputID, outputID string,
) {
	inputChannel := c.registerChannel(inputID)
	outputChannel := c.registerChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyor) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputIDs []string,
	outputID string,
) {
	inputChannels := make([]chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		inputChannels = append(inputChannels, c.registerChannel(id))
	}

	outputChannel := c.registerChannel(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyor) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	inputID string,
	outputIDs []string,
) {
	inputChannel := c.registerChannel(inputID)
	outputChannels := make([]chan string, 0, len(outputIDs))
	for _, id := range outputIDs {
		outputChannels = append(outputChannels, c.registerChannel(id))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyor) Run(parent context.Context) error {
	c.mu.Lock()
	if c.ctx != nil {
		c.mu.Unlock()

		if err := parent.Err(); err != nil {
			return fmt.Errorf("parent context error: %w", err)
		}

		return ErrAlreadyRunning
	}

	c.ctx, c.cancel = context.WithCancel(parent)
	c.mu.Unlock()
	defer c.cancel()

	handlersCopy := append([]handlerFn(nil), c.handlers...)
	errorChannel := make(chan error, len(handlersCopy))

	for _, currentHandler := range handlersCopy {
		c.wg.Add(1)

		go func(h handlerFn) {
			defer c.wg.Done()

			if err := h(c.ctx); err != nil {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}(currentHandler)
	}

	c.wg.Wait()
	close(errorChannel)

	c.mu.Lock()
	for _, channel := range c.chans {
		close(channel)
	}
	c.mu.Unlock()

	for err := range errorChannel {
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return fmt.Errorf("conveyor run failed: %w", err)
		}
	}

	return nil
}

func (c *Conveyor) Send(channelID, value string) error {
	c.mu.Lock()
	channel, exists := c.chans[channelID]
	c.mu.Unlock()

	if !exists {
		return ErrChannelMissing
	}

	select {
	case channel <- value:
		return nil
	default:
	}

	select {
	case channel <- value:
		return nil
	case <-c.getContext().Done():
		return fmt.Errorf("send blocked: %w", c.getContext().Err())
	}
}

func (c *Conveyor) Recv(channelID string) (string, error) {
	c.mu.Lock()
	channel, exists := c.chans[channelID]
	c.mu.Unlock()

	if !exists {
		return "", ErrChannelMissing
	}

	select {
	case value, ok := <-channel:
		if !ok {
			return closedChannelValue, nil
		}
		return value, nil
	default:
	}

	select {
	case value, ok := <-channel:
		if !ok {
			return closedChannelValue, nil
		}
		return value, nil
	case <-c.getContext().Done():
		select {
		case value, ok := <-channel:
			if !ok {
				return closedChannelValue, nil
			}
			return value, nil
		default:
			return "", fmt.Errorf("recv timeout: %w", c.getContext().Err())
		}
	}
}

func (c *Conveyor) getContext() context.Context {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}
