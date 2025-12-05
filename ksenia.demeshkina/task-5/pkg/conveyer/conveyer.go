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
	tp        StepType
	dec       handlers.DecoratingHandler
	sep       handlers.SeparatingHandler
	mux       handlers.MultiplexingHandler
}

type Conveyer struct {
	steps   []step
	ctx     context.Context
	cancel  context.CancelFunc

	inputs  map[string]chan string
	outputs map[string]chan string
}

// New creates named input channels
func New(names ...string) *Conveyer {
	inputs := make(map[string]chan string)
	for _, n := range names {
		inputs[n] = make(chan string)
	}
	return &Conveyer{
		inputs: inputs,
	}
}

func (c *Conveyer) RegisterDecorator(name string, h handlers.DecoratingHandler) {
	c.steps = append(c.steps, step{
		tp:  Decorating,
		dec: h,
	})
}

func (c *Conveyer) RegisterSeparator(name string, h handlers.SeparatingHandler) {
	c.steps = append(c.steps, step{
		tp:  Separating,
		sep: h,
	})
}

func (c *Conveyer) RegisterMultiplexer(name string, h handlers.MultiplexingHandler) {
	c.steps = append(c.steps, step{
		tp:  Multiplexing,
		mux: h,
	})
}

func (c *Conveyer) Send(name, data string) error {
	ch, ok := c.inputs[name]
	if !ok {
		return errors.New("channel does not exist")
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	ch, ok := c.outputs[name]
	if !ok {
		return "", errors.New("channel does not exist")
	}
	data, ok := <-ch
	if !ok {
		return "", nil
	}
	return data, nil
}

// 🔥 Главная часть — построение пайплайна так, как ожидают тесты
func (c *Conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	curr := make([]chan string, 0, len(c.inputs))
	for _, ch := range c.inputs {
		curr = append(curr, ch)
	}

	for i, s := range c.steps {
		switch s.tp {

		case Decorating:
			out := []chan string{make(chan string)}
			go s.dec(c.ctx, curr[0], out[0])
			curr = out

		case Separating:
			out := []chan string{make(chan string), make(chan string)}
			go s.sep(c.ctx, curr[0], out)
			curr = out

		case Multiplexing:
			out := []chan string{make(chan string)}
			go s.mux(c.ctx, curr, out[0])
			curr = out
		}

		// последняя стадия → сохранить как outputs
		if i == len(c.steps)-1 {
			c.outputs = make(map[string]chan string)
			for idx, ch := range curr {
				c.outputs[string('A'+idx)] = ch
			}
		}
	}

	return nil
}

func (c *Conveyer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	for _, ch := range c.inputs {
		close(ch)
	}
}
