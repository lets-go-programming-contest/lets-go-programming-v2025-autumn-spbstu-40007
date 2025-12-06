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

	if existing, ok := c.chans[channelID]; ok {
		return existing
	}

	channel := make(chan string, c.bufSize)
	c.chans[channelID] = channel

	return channel
}

func (c *Conveyor) RegisterDecorator(
	decoratorFn func(context.Context, chan string, chan string) error,
	inID, outID string,
) {
	inputChannel := c.registerChannel(inID)
	outputChannel := c.registerChannel(outID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return decoratorFn(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyor) RegisterMultiplexer(
	multiplexerFn func(context.Context, []chan string, chan string) error,
	inIDs []string,
	outID string,
) {
	inputChannels := make([]chan string, 0, len(inIDs))
	for _, id := range inIDs {
		inputChannels = append(inputChannels, c.registerChannel(id))
	}

	outputChannel := c.registerChannel(outID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return multiplexerFn(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyor) RegisterSeparator(
	separatorFn func(context.Context, chan string, []chan string) error,
	inID string,
	outIDs []string,
) {
	inputChannel := c.registerChannel(inID)

	outputChannels := make([]chan string, 0, len(outIDs))
	for _, id := range outIDs {
		outputChannels = append(outputChannels, c.registerChannel(id))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return separatorFn(ctx, inputChannel, outputChannels)
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

	handlers := append([]handlerFn(nil), c.handlers...)
	errCh := make(chan error, len(handlers))

	for _, handler := range handlers {
		c.wg.Add(1)
		currentHandler := handler
		go func(h handlerFn) {
			defer c.wg.Done()

			err := h(c.ctx)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(currentHandler)
	}

	c.wg.Wait()
	close(errCh)

	c.mu.Lock()
	for _, ch := range c.chans {
		close(ch)
	}
	c.mu.Unlock()

	for err := range errCh {
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
	channel, ok := c.chans[channelID]
	c.mu.Unlock()

	if !ok {
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
	channel, ok := c.chans[channelID]
	c.mu.Unlock()

	if !ok {
		return "", ErrChannelMissing
	}

	select {
	case v, ok := <-channel:
		if !ok {
			return closedChannelValue, nil
		}

		return v, nil
	default:
	}

	select {
	case v, ok := <-channel:
		if !ok {
			return closedChannelValue, nil
		}

		return v, nil
	case <-c.getContext().Done():
		select {
		case v, ok := <-channel:
			if !ok {
				return closedChannelValue, nil
			}

			return v, nil
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
