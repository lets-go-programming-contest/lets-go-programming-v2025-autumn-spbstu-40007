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
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	started   bool
}

func New(bufSize int) *Conveyor {
	return &Conveyor{
		bufSize:  bufSize,
		chMap:    make(map[string]chan string),
		handlers: make([]handlerFn, 0),
	}
}

func (c *Conveyor) createChanIfNeeded(id string) chan string {
	c.syncMutex.Lock()
	defer c.syncMutex.Unlock()

	if ch, ok := c.chMap[id]; ok {

		return ch
	}

	ch := make(chan string, c.bufSize)
	c.chMap[id] = ch

	return ch
}

func (c *Conveyor) getContext() context.Context {
	c.syncMutex.Lock()
	defer c.syncMutex.Unlock()
	if c.ctx == nil {
		return context.Background()
	}

	return c.ctx
}

func (c *Conveyor) ensureStarted() {
	c.syncMutex.Lock()
	defer c.syncMutex.Unlock()

	if c.started || len(c.handlers) == 0 {

		return
	}

	c.started = true
	go c.Run(context.Background())
}

func (c *Conveyor) RegisterDecorator(
	decFn func(ctx context.Context, in chan string, out chan string) error,
	inID, outID string,
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
	c.syncMutex.Lock()
	if c.started {
		c.syncMutex.Unlock()
		<-ctx.Done()

		return ctx.Err()
	}

	c.started = true
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.syncMutex.Unlock()

	errCh := make(chan error, len(c.handlers))

	c.syncMutex.Lock()
	handlers := append([]handlerFn(nil), c.handlers...)
	c.syncMutex.Unlock()

	for _, h := range handlers {
		c.wg.Add(1)
		go func(fn handlerFn) {
			defer c.wg.Done()
			if err := fn(c.ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				c.cancel()
			}
		}(h)
	}

	go func() {
		c.wg.Wait()
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

	c.ensureStarted()

	select {
	case ch <- value:

		return nil
	case <-c.getContext().Done():

		return c.getContext().Err()
	}
}

func (c *Conveyor) Recv(outID string) (string, error) {
	c.syncMutex.Lock()
	ch, found := c.chMap[outID]
	c.syncMutex.Unlock()

	if !found {
		return "", ErrChannelMissing
	}

	c.ensureStarted()

	select {
	case val, ok := <-ch:
		if !ok {

			return "undefined", nil
		}
		return val, nil
	case <-c.getContext().Done():

		return "", c.getContext().Err()
	}
}
