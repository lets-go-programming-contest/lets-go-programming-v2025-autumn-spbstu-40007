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
	tp  StepType
	dec handlers.DecoratingHandler
	sep handlers.SeparatingHandler
	mux handlers.MultiplexingHandler
}

type Conveyer struct {
	steps   []step
	ctx     context.Context
	cancel  context.CancelFunc

	input  chan string
	outputs map[string]chan string
}

func New(_ int) *Conveyer {
	return &Conveyer{
		input: make(chan string),
	}
}

// тесты ожидают только один параметр — data
func (c *Conveyer) Send(data string) error {
	if c.input == nil {
		return errors.New("no input")
	}
	c.input <- data
	return nil
}

// тесты вызывают Recv() без параметров
func (c *Conveyer) Recv() (string, error) {
	if c.outputs == nil {
		return "", errors.New("pipeline not started")
	}
	var out chan string
	for _, o := range c.outputs {
		out = o
		break
	}
	data, ok := <-out
	if !ok {
		return "", nil
	}
	return data, nil
}

func (c *Conveyer) RegisterDecorator(_ string, h handlers.DecoratingHandler) {
	c.steps = append(c.steps, step{tp: Decorating, dec: h})
}

func (c *Conveyer) RegisterSeparator(_ string, h handlers.SeparatingHandler) {
	c.steps = append(c.steps, step{tp: Separating, sep: h})
}

func (c *Conveyer) RegisterMultiplexer(_ string, h handlers.MultiplexingHandler) {
	c.steps = append(c.steps, step{tp: Multiplexing, mux: h})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	curr := []chan string{c.input}

	for i, s := range c.steps {
		switch s.tp {

		case Decorating:
			out := make(chan string)
			go s.dec(c.ctx, curr[0], out)
			curr = []chan string{out}

		case Separating:
			outs := []chan string{make(chan string), make(chan string)}
			go s.sep(c.ctx, curr[0], outs)
			curr = outs

		case Multiplexing:
			out := make(chan string)
			go s.mux(c.ctx, curr, out)
			curr = []chan string{out}
		}

		if i == len(c.steps)-1 {
			c.outputs = make(map[string]chan string)
			if len(curr) == 1 {
				c.outputs["data"] = curr[0]
			} else if len(curr) == 2 {
				c.outputs["ok"] = curr[0]
				c.outputs["err"] = curr[1]
			}
		}
	}

	return nil
}

func (c *Conveyer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	if c.input != nil {
		close(c.input)
	}
}
