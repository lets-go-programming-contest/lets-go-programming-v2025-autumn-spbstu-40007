package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrClosed       = errors.New("channel closed")
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type conveyer interface {
	RegisterDecorator(fn DecoratorFunc, input string, output string)
	RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string)
	RegisterSeparator(fn SeparatorFunc, input string, outputs []string)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type HandlerEntry struct {
	Name    string
	RunFunc func(ctx context.Context, wg *sync.WaitGroup, channels map[string]chan string) error
}

type Conveyer struct {
	size int

	mu       sync.RWMutex
	channels map[string]chan string

	handlers []HandlerEntry
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]HandlerEntry, 0),
	}
}

func (c *Conveyer) getOrCreateChannel(name string) (chan string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch, false
	}

	newCh := make(chan string, c.size)
	c.channels[name] = newCh
	return newCh, true
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[name]
	if !exists {
		return nil, ErrChanNotFound
	}
	return ch, nil
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()
			close(ch)
		}()
	}
}

func (c *Conveyer) RegisterDecorator(fn DecoratorFunc, input string, output string) {
	c.getOrCreateChannel(input)
	c.getOrCreateChannel(output)

	entry := HandlerEntry{
		Name: fmt.Sprintf("Decorator(%s -> %s)", input, output),
		RunFunc: func(ctx context.Context, wg *sync.WaitGroup, channels map[string]chan string) error {
			defer wg.Done()
			return fn(ctx, channels[input], channels[output])
		},
	}
	c.handlers = append(c.handlers, entry)
}

func (c *Conveyer) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	for _, name := range inputs {
		c.getOrCreateChannel(name)
	}
	c.getOrCreateChannel(output)

	entry := HandlerEntry{
		Name: fmt.Sprintf("Multiplexer(%v -> %s)", inputs, output),
		RunFunc: func(ctx context.Context, wg *sync.WaitGroup, channels map[string]chan string) error {
			defer wg.Done()
			inputChans := make([]chan string, len(inputs))
			for i, name := range inputs {
				inputChans[i] = channels[name]
			}
			return fn(ctx, inputChans, channels[output])
		},
	}
	c.handlers = append(c.handlers, entry)
}

func (c *Conveyer) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	c.getOrCreateChannel(input)
	for _, name := range outputs {
		c.getOrCreateChannel(name)
	}

	entry := HandlerEntry{
		Name: fmt.Sprintf("Separator(%s -> %v)", input, outputs),
		RunFunc: func(ctx context.Context, wg *sync.WaitGroup, channels map[string]chan string) error {
			defer wg.Done()
			outputChans := make([]chan string, len(outputs))
			for i, name := range outputs {
				outputChans[i] = channels[name]
			}
			return fn(ctx, channels[input], outputChans)
		},
	}
	c.handlers = append(c.handlers, entry)
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.handlers))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer c.closeAllChannels()

	c.mu.RLock()
	currentChannels := c.channels
	c.mu.RUnlock()

	for _, handler := range c.handlers {
		wg.Add(1)
		go func(entry HandlerEntry) {
			defer wg.Done()
			if err := entry.RunFunc(ctx, &wg, currentChannels); err != nil {
				errCh <- fmt.Errorf("%s failed: %w", entry.Name, err)
				cancel()
			}
		}(handler)
	}

	go func() {
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		if err, ok := <-errCh; ok {
			return err
		}

		return ctx.Err()

	case err, ok := <-errCh:
		if ok {
			return err
		}

		return nil
	}
}

func (c *Conveyer) Send(input string, data string) error {
	ch, err := c.getChannel(input)
	if err != nil {
		return ErrChanNotFound
	}
	select {
	case ch <- data:
		return nil
	default:
		return fmt.Errorf("send to channel %s failed: channel blocked or closed", input)
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, err := c.getChannel(output)
	if err != nil {
		return "undefined", ErrChanNotFound
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
}
