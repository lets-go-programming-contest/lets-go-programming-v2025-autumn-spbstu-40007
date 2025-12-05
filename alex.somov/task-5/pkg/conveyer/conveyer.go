package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type decoratorDesc struct {
	fn     func(ctx context.Context, input chan string, output chan string) error
	input  string
	output string
}

type multiplexerDesc struct {
	fn     func(ctx context.Context, inputs []chan string, output chan string) error
	inputs []string
	output string
}

type separatorDesc struct {
	fn      func(ctx context.Context, input chan string, outputs []chan string) error
	input   string
	outputs []string
}

type Conveyer struct {
	chans map[string]chan string

	decorators   []decoratorDesc
	multiplexers []multiplexerDesc
	separators   []separatorDesc

	bufSize int
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func New(size int) *Conveyer {
	return &Conveyer{
		chans:        make(map[string]chan string),
		decorators:   make([]decoratorDesc, 0),
		multiplexers: make([]multiplexerDesc, 0),
		separators:   make([]separatorDesc, 0),
		bufSize:      size,
	}
}

func (c *Conveyer) getOrCreate(name string) chan string {
	ch, ok := c.chans[name]
	if !ok {
		ch = make(chan string, c.bufSize)
		c.chans[name] = ch
	}

	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreate(input)
	c.getOrCreate(output)

	c.decorators = append(c.decorators, decoratorDesc{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, id := range inputs {
		c.getOrCreate(id)
	}
	c.getOrCreate(output)

	c.multiplexers = append(c.multiplexers, multiplexerDesc{
		fn:     fn,
		inputs: append([]string(nil), inputs...),
		output: output,
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreate(input)
	for _, id := range outputs {
		c.getOrCreate(id)
	}

	c.separators = append(c.separators, separatorDesc{
		fn:      fn,
		input:   input,
		outputs: append([]string(nil), outputs...),
	})
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	ch, ok := c.chans[input]
	c.mu.Unlock()

	if !ok {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	ch, ok := c.chans[output]
	c.mu.Unlock()

	if !ok {
		return "", ErrChanNotFound
	}

	v, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return v, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	c.mu.Lock()

	for _, d := range c.decorators {
		in := c.chans[d.input]
		out := c.chans[d.output]

		c.wg.Add(1)
		go func(desc decoratorDesc, in, out chan string) {
			defer c.wg.Done()
			if err := desc.fn(ctx, in, out); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(d, in, out)
	}

	for _, m := range c.multiplexers {
		var ins []chan string
		for _, id := range m.inputs {
			ins = append(ins, c.chans[id])
		}
		out := c.chans[m.output]

		c.wg.Add(1)
		go func(desc multiplexerDesc, ins []chan string, out chan string) {
			defer c.wg.Done()
			if err := desc.fn(ctx, ins, out); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(m, ins, out)
	}

	for _, s := range c.separators {
		in := c.chans[s.input]
		var outs []chan string
		for _, id := range s.outputs {
			outs = append(outs, c.chans[id])
		}

		c.wg.Add(1)
		go func(desc separatorDesc, in chan string, outs []chan string) {
			defer c.wg.Done()
			if err := desc.fn(ctx, in, outs); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(s, in, outs)
	}

	c.mu.Unlock()

	select {
	case err := <-errCh:
		cancel()
		c.wg.Wait()
		return err
	case <-ctx.Done():
		c.wg.Wait()
		return nil
	}
}
