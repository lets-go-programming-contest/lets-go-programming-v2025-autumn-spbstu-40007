package conveyer

import (
	"context"
	"errors"
	"sync"
)

const (
	errMsgChanNotFound = "chan not found"
	valUndefined       = "undefined"
)

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	tasks    []func(context.Context) error
	size     int
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		tasks:    make([]func(context.Context) error, 0),
		size:     size,
	}
}

func (c *Conveyer) getOrInitChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(fn func(context.Context, chan string, chan string) error, inputName, outputName string) {
	inputCh := c.getOrInitChannel(inputName)
	outputCh := c.getOrInitChannel(outputName)

	task := func(ctx context.Context) error {
		defer close(outputCh)
		return fn(ctx, inputCh, outputCh)
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inputNames []string, outputName string) {
	var inputs []chan string
	for _, name := range inputNames {
		inputs = append(inputs, c.getOrInitChannel(name))
	}
	outputCh := c.getOrInitChannel(outputName)

	task := func(ctx context.Context) error {
		defer close(outputCh)
		return fn(ctx, inputs, outputCh)
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterSeparator(fn func(context.Context, chan string, []chan string) error, inputName string, outputNames []string) {
	inputCh := c.getOrInitChannel(inputName)
	var outputs []chan string
	for _, name := range outputNames {
		outputs = append(outputs, c.getOrInitChannel(name))
	}

	task := func(ctx context.Context) error {
		defer func() {
			for _, ch := range outputs {
				close(ch)
			}
		}()
		return fn(ctx, inputCh, outputs)
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, len(c.tasks))

	for _, task := range c.tasks {
		wg.Add(1)
		go func(t func(context.Context) error) {
			defer wg.Done()
			if err := t(ctx); err != nil {
				select {
				case errCh <- err:
					cancel()
				default:
				}
			}
		}(task)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-done:
		return nil
	}
}

func (c *Conveyer) Send(name string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[name]
	c.mu.RUnlock()

	if !ok {
		return errors.New(errMsgChanNotFound)
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[name]
	c.mu.RUnlock()

	if !ok {
		return "", errors.New(errMsgChanNotFound)
	}

	val, isOpen := <-ch
	if !isOpen {
		return valUndefined, nil
	}
	return val, nil
}
