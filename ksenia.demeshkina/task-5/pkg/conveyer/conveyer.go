package conveyer

import (
	"context"
	"errors"

	"github.com/ksuah/task-5/pkg/handlers"
)

type StepType int

const (
	Decorating StepType = iota
	Separating
	Multiplexing
)

type step struct {
	tp   StepType
	name string

	in     string
	out    string
	outOk  string
	outErr string

	dec handlers.DecoratingHandler
	sep handlers.SeparatingHandler
	mux handlers.MultiplexingHandler
}

type Conveyer struct {
	ctx    context.Context
	cancel context.CancelFunc

	steps []step

	chans map[string]chan string
}

func New(_ int) *Conveyer {
	return &Conveyer{
		chans: make(map[string]chan string),
	}
}

// Создаём канал при первом обращении
func (c *Conveyer) getChan(name string) chan string {
	ch, ok := c.chans[name]
	if !ok {
		ch = make(chan string)
		c.chans[name] = ch
	}
	return ch
}

func (c *Conveyer) Send(name string, data string) error {
	ch := c.getChan(name)
	ch <- data
	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	ch, ok := c.chans[name]
	if !ok {
		return "", errors.New("channel does not exist")
	}
	data, ok := <-ch
	if !ok {
		return "", nil
	}
	return data, nil
}

func (c *Conveyer) RegisterDecorator(name, in, out string, h handlers.DecoratingHandler) {
	c.steps = append(c.steps, step{
		name: name,
		tp:   Decorating,
		in:   in,
		out:  out,
		dec:  h,
	})
}

func (c *Conveyer) RegisterSeparator(name, in, outOk, outErr string, h handlers.SeparatingHandler) {
	c.steps = append(c.steps, step{
		name:   name,
		tp:     Separating,
		in:     in,
		outOk:  outOk,
		outErr: outErr,
		sep:    h,
	})
}

func (c *Conveyer) RegisterMultiplexer(name, in1, in2, out string, h handlers.MultiplexingHandler) {
	c.steps = append(c.steps, step{
		name: name,
		tp:   Multiplexing,
		in:   in1 + "," + in2,
		out:  out,
		mux:  h,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	for _, s := range c.steps {
		switch s.tp {

		case Decorating:
			in := c.getChan(s.in)
			out := c.getChan(s.out)
			go s.dec(c.ctx, in, out)

		case Separating:
			in := c.getChan(s.in)
			ok := c.getChan(s.outOk)
			errCh := c.getChan(s.outErr)
			go s.sep(c.ctx, in, []chan string{ok, errCh})

		case Multiplexing:
			parts := []string{}
			for _, p := range []rune(s.in) {
				if p != ',' {
					parts = append(parts, string(p))
				}
			}
			var inputs []chan string
			for _, nm := range parts {
				inputs = append(inputs, c.getChan(nm))
			}
			out := c.getChan(s.out)
			go s.mux(c.ctx, inputs, out)
		}
	}
	return nil
}

func (c *Conveyer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	for _, ch := range c.chans {
		close(ch)
	}
}
