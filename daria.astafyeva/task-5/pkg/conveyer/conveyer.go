package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrChannelMissing = errors.New("chan not found")

type handlerFn func(ctx context.Context) error

type Conveyor struct {
	bufSize   int
	chMap     map[string]chan string
	handlers  []handlerFn
	syncMutex sync.Mutex
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
}

func New(bufSize int) *Conveyor {
	return &Conveyor{
		bufSize:  bufSize,
		chMap:    make(map[string]chan string),
		handlers: make([]handlerFn, 0),
		running:  false,
	}
}

func (c *Conveyor) createChanIfNeeded(id string) chan string {
	c.syncMutex.Lock()
	defer c.syncMutex.Unlock()

	if ch, found := c.chMap[id]; found {
		return ch
	}

	newCh := make(chan string, c.bufSize)
	c.chMap[id] = newCh
	return newCh
}

func (c *Conveyor) RegisterDecorator(
	decFn func(ctx context.Context, in chan string, out chan string) error,
	inID string,
	outID string,
) {
	inCh := c.createChanIfNeeded(inID)
	outCh := c.createChanIfNeeded(outID)

	c.syncMutex.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return decFn(ctx, inCh, outCh)
	})
	c.syncMutex.Unlock()
}

func (c *Conveyor) RegisterMultiplexer(
	muxFn func(ctx context.Context, ins []chan string, out chan string) error,
	inIDs []string,
	outID string,
) {
	inChs := make([]chan string, 0, len(inIDs))
	for _, id := range inIDs {
		inChs = append(inChs, c.createChanIfNeeded(id))
	}
	outCh := c.createChanIfNeeded(outID)

	c.syncMutex.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return muxFn(ctx, inChs, outCh)
	})
	c.syncMutex.Unlock()
}

func (c *Conveyor) RegisterSeparator(
	sepFn func(ctx context.Context, in chan string, outs []chan string) error,
	inID string,
	outIDs []string,
) {
	inCh := c.createChanIfNeeded(inID)
	outChs := make([]chan string, 0, len(outIDs))
	for _, id := range outIDs {
		outChs = append(outChs, c.createChanIfNeeded(id))
	}

	c.syncMutex.Lock()
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return sepFn(ctx, inCh, outChs)
	})
	c.syncMutex.Unlock()
}

func (c *Conveyor) Run(ctx context.Context) error {
	if c.running {
		return errors.New("conveyor already running")
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()

	c.running = true
	defer func() { c.running = false }()

	var wg sync.WaitGroup
	errCh := make(chan error, len(c.handlers))

	c.syncMutex.Lock()
	for _, h := range c.handlers {
		wg.Add(1)
		go func(fn handlerFn) {
			defer wg.Done()
			if err := fn(c.ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				c.cancel()
			}
		}(h)
	}
	c.syncMutex.Unlock()

	go func() {
		wg.Wait()
		close(errCh)

		c.syncMutex.Lock()
		for _, ch := range c.chMap {
			close(ch)
		}
		c.syncMutex.Unlock()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("conveyor run failed: %w", err)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conveyor) Send(inID string, value string) error {
	c.syncMutex.Lock()
	ch, found := c.chMap[inID]
	c.syncMutex.Unlock()

	if !found {
		return ErrChannelMissing
	}

	select {
	case ch <- value:
		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
}

func (c *Conveyor) Recv(outID string) (string, error) {
	c.syncMutex.Lock()
	ch, found := c.chMap[outID]
	c.syncMutex.Unlock()

	if !found {
		return "", ErrChannelMissing
	}

	select {
	case val, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return val, nil
	case <-c.ctx.Done():
		return "", c.ctx.Err()
	}
}
