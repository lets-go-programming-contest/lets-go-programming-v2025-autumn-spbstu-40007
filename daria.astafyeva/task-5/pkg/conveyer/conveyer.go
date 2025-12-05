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
}

func New(size int) *Conveyor {
	return &Conveyor{
		bufSize: size,
		chans:   make(map[string]chan string),
	}
}

func (c *Conveyor) registerChan(id string) chan string {
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
	inCh := c.registerChan(inID)
	outCh := c.registerChan(outID)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *Conveyor) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inIDs []string, outID string) {
	ins := make([]chan string, len(inIDs))
	for i, id := range inIDs {
		ins[i] = c.registerChan(id)
	}
	out := c.registerChan(outID)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (c *Conveyor) RegisterSeparator(fn func(context.Context, chan string, []chan string) error, inID string, outIDs []string) {
	in := c.registerChan(inID)
	outs := make([]chan string, len(outIDs))
	for i, id := range outIDs {
		outs[i] = c.registerChan(id)
	}
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outs)
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
	c.mu.Unlock()

	handlers := append([]handlerFn(nil), c.handlers...)

	errCh := make(chan error, len(handlers))

	for _, h := range handlers {
		c.wg.Add(1)
		go func(fn handlerFn) {
			defer c.wg.Done()
			if err := fn(c.ctx); err != nil {
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
		for _, ch := range c.chans {
			close(ch)
		}
		c.mu.Unlock()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			c.cancel()
			return fmt.Errorf("conveyor run failed: %w", err)
		}

		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
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
	case <-c.getCtx().Done():
		return c.getCtx().Err()
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
			return "undefined", nil
		}
		return v, nil
	case <-c.getCtx().Done():
		return "", c.getCtx().Err()
	}
}

func (c *Conveyor) getCtx() context.Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}
