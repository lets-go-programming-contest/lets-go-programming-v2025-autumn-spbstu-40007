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

func New(bufSize int) *Conveyor {
	return &Conveyor{
		bufSize:  bufSize,
		chans:    make(map[string]chan string),
		handlers: make([]handlerFn, 0),
	}
}

func (c *Conveyor) registerChan(channelID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.chans[channelID]; ok {
		return ch
	}

	ch := make(chan string, c.bufSize)
	c.chans[channelID] = ch

	return ch
}

func (c *Conveyor) RegisterDecorator(
	decoratorFn func(context.Context, chan string, chan string) error,
	inputID, outputID string,
) {
	input := c.registerChan(inputID)
	output := c.registerChan(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return decoratorFn(ctx, input, output)
	})
}

func (c *Conveyor) RegisterMultiplexer(
	multiplexerFn func(context.Context, []chan string, chan string) error,
	inputIDs []string,
	outputID string,
) {
	inputs := make([]chan string, 0, len(inputIDs))
	for _, id := range inputIDs {
		inputs = append(inputs, c.registerChan(id))
	}
	output := c.registerChan(outputID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return multiplexerFn(ctx, inputs, output)
	})
}

func (c *Conveyor) RegisterSeparator(
	separatorFn func(context.Context, chan string, []chan string) error,
	inputID string,
	outputIDs []string,
) {
	input := c.registerChan(inputID)
	outputs := make([]chan string, 0, len(outputIDs))
	for _, id := range outputIDs {
		outputs = append(outputs, c.registerChan(id))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return separatorFn(ctx, input, outputs)
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
		go func(h handlerFn) {
			defer c.wg.Done()

			if err := h(c.ctx); err != nil {
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
	ch, ok := c.chans[channelID]
	c.mu.Unlock()
	if !ok {
		return ErrChannelMissing
	}

	select {
	case ch <- value:
		return nil
	default:
	}

	select {
	case ch <- value:
		return nil
	case <-c.getContext().Done():
		return fmt.Errorf("send blocked: %w", c.getContext().Err())
	}
}

func (c *Conveyor) Recv(channelID string) (string, error) {
	c.mu.Lock()
	ch, ok := c.chans[channelID]
	c.mu.Unlock()
	if !ok {
		return "", ErrChannelMissing
	}

	select {
	case v, ok := <-ch:
		if !ok {
			return closedChannelValue, nil
		}
		return v, nil
	default:
	}

	select {
	case v, ok := <-ch:
		if !ok {
			return closedChannelValue, nil
		}
		return v, nil
	case <-c.getContext().Done():
		select {
		case v, ok := <-ch:
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
