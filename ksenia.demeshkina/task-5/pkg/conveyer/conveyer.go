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
	done       chan struct{}
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:         sync.Mutex{},
		channels:   make(map[string]chan string),
		configs:    []handlerConfig{},
		bufferSize: size,
		wg:         sync.WaitGroup{},
		done:       make(chan struct{}),
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
	handlerFunc func(context.Context, chan string, chan string) error,
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
		fn:      handlerFunc,
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
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
		fn:      handlerFunc,
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
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
		fn:      handlerFunc,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 1)
	c.done = make(chan struct{})

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
				fn, ok := conf.fn.(func(context.Context, chan string, chan string) error)
				if !ok {
					err = errors.New("invalid decorator function type")
				} else {
					err = fn(ctx, inChs[0], outChs[0])
				}

			case typeMultiplexer:
				fn, ok := conf.fn.(func(context.Context, []chan string, chan string) error)
				if !ok {
					err = errors.New("invalid multiplexer function type")
				} else {
					err = fn(ctx, inChs, outChs[0])
				}

			case typeSeparator:
				fn, ok := conf.fn.(func(context.Context, chan string, []chan string) error)
				if !ok {
					err = errors.New("invalid separator function type")
				} else {
					err = fn(ctx, inChs[0], outChs)
				}
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

	go func() {
		c.wg.Wait()
		close(c.done)
	}()

	var resultErr error

	select {
	case err := <-errChan:
		resultErr = err
	case <-ctx.Done():
		resultErr = nil
	case <-c.done:
		resultErr = nil
	}

	<-c.done

	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
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
