package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrChannelMissing = errors.New("chan not found")

type handlerFn func(context.Context) error

type Conveyor struct {
	bufSize  int
	chans    map[string]chan string
	handlers []handlerFn
	mu       sync.Mutex
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	closed   bool
}

func New(size int) *Conveyor {
	return &Conveyor{
		bufSize: size,
		chans:   make(map[string]chan string),
	}
}

func (c *Conveyor) ensureChan(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.chans[id]; ok {
		return ch
	}
	ch := make(chan string, c.bufSize)
	c.chans[id] = ch
	return ch
}

func (c *Conveyor) getChan(id string) (chan string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch, ok := c.chans[id]
	if !ok {
		return nil, ErrChannelMissing
	}
	return ch, nil
}

func (c *Conveyor) RegisterDecorator(fn func(context.Context, chan string, chan string) error, inID, outID string) {
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, c.ensureChan(inID), c.ensureChan(outID))
	})
}

func (c *Conveyor) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inIDs []string, outID string) {
	ins := make([]chan string, len(inIDs))
	for i, id := range inIDs {
		ins[i] = c.ensureChan(id)
	}
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, ins, c.ensureChan(outID))
	})
}

func (c *Conveyor) RegisterSeparator(fn func(context.Context, chan string, []chan string) error, inID string, outIDs []string) {
	outs := make([]chan string, len(outIDs))
	for i, id := range outIDs {
		outs[i] = c.ensureChan(id)
	}
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, c.ensureChan(inID), outs)
	})
}

func (c *Conveyor) Run(parent context.Context) error {
	c.mu.Lock()
	if c.ctx != nil {
		c.mu.Unlock()
		<-parent.Done()
		return parent.Err()
	}
	c.ctx, c.cancel = context.WithCancel(parent)
	handlers := append([]handlerFn{}, c.handlers...)
	c.mu.Unlock()

	errCh := make(chan error, len(handlers))

	for _, h := range handlers {
		c.wg.Add(1)
		go func(h handlerFn) {
			defer c.wg.Done()
			if err := h(c.ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(h)
	}

	go func() {
		c.wg.Wait()
		close(errCh)

		c.mu.Lock()
		defer c.mu.Unlock()
		if c.closed {
			return
		}
		c.closed = true
		for _, ch := range c.chans {
			close(ch)
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			c.cancel()
			return fmt.Errorf("conveyor run failed: %w", err)
		}
		return nil
	case <-parent.Done():
		c.cancel()
		return parent.Err()
	}
}

func (c *Conveyor) Send(id, val string) error {
	ch, err := c.getChan(id)
	if err != nil {
		return err
	}
	select {
	case ch <- val:
		return nil
	case <-c.ctxOrBg().Done():
		return c.ctxOrBg().Err()
	}
}

func (c *Conveyor) Recv(id string) (string, error) {
	ch, err := c.getChan(id)
	if err != nil {
		return "", err
	}
	select {
	case v, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return v, nil
	case <-c.ctxOrBg().Done():
		return "", c.ctxOrBg().Err()
	}
}

func (c *Conveyor) ctxOrBg() context.Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}
