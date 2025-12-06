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
	bufSize        int
	chans          map[string]chan string
	handlers       []handlerFn
	mu             sync.Mutex
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	canceledCtx    context.Context
	canceledCancel context.CancelFunc
}

func New(size int) *Conveyor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Conveyor{
		bufSize:        size,
		chans:          make(map[string]chan string),
		canceledCtx:    ctx,
		canceledCancel: cancel,
	}
}

func (c *Conveyor) registerChan(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.chans[id]; exists {
		return ch
	}

	ch := make(chan string, c.bufSize)
	c.chans[id] = ch
	return ch
}

func (c *Conveyor) RegisterDecorator(fn func(context.Context, chan string, chan string) error, inID, outID string) {
	inCh := c.registerChan(inID)
	outCh := c.registerChan(outID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *Conveyor) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inIDs []string, outID string) {
	inputs := make([]chan string, len(inIDs))
	for i, id := range inIDs {
		inputs[i] = c.registerChan(id)
	}
	output := c.registerChan(outID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputs, output)
	})
}

func (c *Conveyor) RegisterSeparator(fn func(context.Context, chan string, []chan string) error, inID string, outIDs []string) {
	input := c.registerChan(inID)
	outputs := make([]chan string, len(outIDs))
	for i, id := range outIDs {
		outputs[i] = c.registerChan(id)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, input, outputs)
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
	errCh := make(chan error, len(handlersCopy))

	for _, handler := range handlersCopy {
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

func (c *Conveyor) Send(id, val string) error {
	c.mu.Lock()
	ch, ok := c.chans[id]
	c.mu.Unlock()
	if !ok {
		return ErrChannelMissing
	}

	select {
	case ch <- val:
		return nil
	default:
	}

	select {
	case ch <- val:
		return nil
	case <-c.getCtx().Done():
		return fmt.Errorf("send failed: %w", c.getCtx().Err())
	}
}

func (c *Conveyor) Recv(id string) (string, error) {
	c.mu.Lock()
	ch, ok := c.chans[id]
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
	case <-c.getCtx().Done():
		select {
		case v, ok := <-ch:
			if !ok {
				return closedChannelValue, nil
			}
			return v, nil
		default:
			return "", fmt.Errorf("recv timeout: %w", c.getCtx().Err())
		}
	}
}

func (c *Conveyor) getCtx() context.Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ctx == nil {
		return c.canceledCtx
	}
	return c.ctx
}
