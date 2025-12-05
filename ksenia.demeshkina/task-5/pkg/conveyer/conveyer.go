package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

const (
	typeDecorator   = 1
	typeMultiplexer = 2
	typeSeparator   = 3
)

type handlerConfig struct {
	hType   int
	inputs  []string
	outputs []string
	fn      interface{}
}

type Conveyer struct {
	mu         sync.Mutex
	channels   map[string]chan string
	configs    []handlerConfig
	bufferSize int
	wg         sync.WaitGroup
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels:   make(map[string]chan string),
		configs:    make([]handlerConfig, 0),
		bufferSize: size,
	}
}

func (c *Conveyer) getOrCreateChan(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChan(input)
	c.getOrCreateChan(output)

	c.configs = append(c.configs, handlerConfig{
		hType:   typeDecorator,
		inputs:  []string{input},
		outputs: []string{output},
		fn:      fn,
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, in := range inputs {
		c.getOrCreateChan(in)
	}
	c.getOrCreateChan(output)

	c.configs = append(c.configs, handlerConfig{
		hType:   typeMultiplexer,
		inputs:  inputs,
		outputs: []string{output},
		fn:      fn,
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChan(input)
	for _, out := range outputs {
		c.getOrCreateChan(out)
	}

	c.configs = append(c.configs, handlerConfig{
		hType:   typeSeparator,
		inputs:  []string{input},
		outputs: outputs,
		fn:      fn,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 1)

	for _, cfg := range c.configs {
		c.wg.Add(1)

		go func(conf handlerConfig) {
			defer c.wg.Done()

			c.mu.Lock()
			var inChs []chan string
			for _, name := range conf.inputs {
				inChs = append(inChs, c.channels[name])
			}
			var outChs []chan string
			for _, name := range conf.outputs {
				outChs = append(outChs, c.channels[name])
			}
			c.mu.Unlock()

			var err error
			switch conf.hType {
			case typeDecorator:
				fn := conf.fn.(func(context.Context, chan string, chan string) error)
				err = fn(ctx, inChs[0], outChs[0])
			case typeMultiplexer:
				fn := conf.fn.(func(context.Context, []chan string, chan string) error)
				err = fn(ctx, inChs, outChs[0])
			case typeSeparator:
				fn := conf.fn.(func(context.Context, chan string, []chan string) error)
				err = fn(ctx, inChs[0], outChs)
			}

			if err != nil {
				select {
				case errChan <- err:
					cancel()
				default:
				}
			}
		}(cfg)
	}

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	var resultErr error

	select {
	case resultErr = <-errChan:
	case <-ctx.Done():
		resultErr = ctx.Err()
	case <-done:
		resultErr = nil
	}

	if resultErr != nil {
		<-done
	}

	c.mu.Lock()
	for _, ch := range c.channels {
		select {
		case <-ch:
		default:
			close(ch)
		}
	}
	c.mu.Unlock()

	return resultErr
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	ch, ok := c.channels[input]
	c.mu.Unlock()

	if !ok {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	ch, ok := c.channels[output]
	c.mu.Unlock()

	if !ok {
		return "", ErrChanNotFound
	}

	data, open := <-ch
	if !open {
		return "undefined", nil
	}

	return data, nil
}