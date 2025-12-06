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
	cancel   context.CancelFunc
	running  bool
}

func New(bufferSize int) *Conveyor {
	return &Conveyor{
		bufSize:  bufferSize,
		chans:    make(map[string]chan string),
		handlers: make([]handlerFn, 0),
		running:  false,
	}
}

func (c *Conveyor) registerChannel(channelID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.chans[channelID]; ok {
		return ch
	}

	channel := make(chan string, c.bufSize)
	c.chans[channelID] = channel
	return channel
}

func (c *Conveyor) RegisterDecorator(
	decoratorFn func(context.Context, chan string, chan string) error,
	inID, outID string,
) {
	inCh := c.registerChannel(inID)
	outCh := c.registerChannel(outID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return decoratorFn(ctx, inCh, outCh)
	})
}

func (c *Conveyor) RegisterMultiplexer(
	multiplexerFn func(context.Context, []chan string, chan string) error,
	inIDs []string,
	outID string,
) {
	inChs := make([]chan string, 0, len(inIDs))
	for _, id := range inIDs {
		inChs = append(inChs, c.registerChannel(id))
	}
	outCh := c.registerChannel(outID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return multiplexerFn(ctx, inChs, outCh)
	})
}

func (c *Conveyor) RegisterSeparator(
	separatorFn func(context.Context, chan string, []chan string) error,
	inID string,
	outIDs []string,
) {
	inCh := c.registerChannel(inID)
	outChs := make([]chan string, 0, len(outIDs))
	for _, id := range outIDs {
		outChs = append(outChs, c.registerChannel(id))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return separatorFn(ctx, inCh, outChs)
	})
}

func (c *Conveyor) Run(parent context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		if err := parent.Err(); err != nil {
			return fmt.Errorf("parent context error: %w", err)
		}
		return ErrAlreadyRunning
	}

	ctx, cancel := context.WithCancel(parent)
	c.cancel = cancel
	c.running = true
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		c.cancel = nil
		c.running = false
		c.mu.Unlock()
	}()

	handlers := append([]handlerFn(nil), c.handlers...)
	errCh := make(chan error, len(handlers))

	for _, handler := range handlers {
		c.wg.Add(1)
		go func(h handlerFn) {
			defer c.wg.Done()
			if err := h(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(handler)
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

	if c.cancel == nil {
		return fmt.Errorf("send blocked: conveyor not running")
	}

	select {
	case channel <- value:
		return nil
	case <-context.Background().Done():
	}
	return fmt.Errorf("send blocked: channel full")
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

	if c.cancel == nil {
		return "", fmt.Errorf("recv timeout: conveyor not running")
	}

	select {
	case v, ok := <-channel:
		if !ok {
			return closedChannelValue, nil
		}
		return v, nil
	case <-context.Background().Done():
		return "", fmt.Errorf("recv timeout")
	}
}
