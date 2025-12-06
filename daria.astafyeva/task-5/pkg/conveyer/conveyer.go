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
	}
}

func (c *Conveyor) registerChannel(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.chans[id]; ok {
		return ch
	}

	ch := make(chan string, c.bufSize)
	c.chans[id] = ch

	return ch
}

func (c *Conveyor) RegisterDecorator(fn func(context.Context, chan string, chan string) error, inID, outID string) {
	in := c.registerChannel(inID)
	out := c.registerChannel(outID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyor) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inIDs []string, outID string) {
	ins := make([]chan string, 0, len(inIDs))
	for _, id := range inIDs {
		ins = append(ins, c.registerChannel(id))
	}
	out := c.registerChannel(outID)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (c *Conveyor) RegisterSeparator(fn func(context.Context, chan string, []chan string) error, inID string, outIDs []string) {
	in := c.registerChannel(inID)
	outs := make([]chan string, 0, len(outIDs))
	for _, id := range outIDs {
		outs = append(outs, c.registerChannel(id))
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outs)
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

	for _, h := range handlers {
		c.wg.Add(1)
		go func(handler handlerFn) {
			defer c.wg.Done()

			if err := handler(c.ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(h)
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

func (c *Conveyor) Send(id, value string) error {
	c.mu.Lock()
	ch, ok := c.chans[id]
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
