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

type Step struct {
	Type      StepType
	Decorator handlers.DecoratingHandler
	Separator handlers.SeparatingHandler
	Multiplex handlers.MultiplexingHandler
}

type Conveyer struct {
	steps []Step

	inputs  []chan string
	outputs []chan string
}

func New() *Conveyer {
	return &Conveyer{}
}

func (c *Conveyer) Send(data string, chIdx int) error {
	if chIdx < 0 || chIdx >= len(c.inputs) {
		return errors.New("channel does not exist")
	}
	c.inputs[chIdx] <- data
	return nil
}

func (c *Conveyer) Recv(chIdx int) (string, error) {
	if chIdx < 0 || chIdx >= len(c.outputs) {
		return "", errors.New("channel does not exist")
	}
	data, ok := <-c.outputs[chIdx]
	if !ok {
		return "", nil
	}
	return data, nil
}

func (c *Conveyer) AddStep(step Step) {
	c.steps = append(c.steps, step)
}

func (c *Conveyer) Run(ctx context.Context) error {
	if len(c.steps) == 0 {
		return nil
	}

	c.inputs = []chan string{make(chan string)}
	curr := c.inputs

	for _, step := range c.steps {
		switch step.Type {

		case Decorating:
			out := []chan string{make(chan string)}
			go step.Decorator(ctx, curr[0], out[0])
			curr = out

		case Separating:
			out := []chan string{make(chan string), make(chan string)}
			go step.Separator(ctx, curr[0], out)
			curr = out

		case Multiplexing:
			out := []chan string{make(chan string)}
			go step.Multiplex(ctx, curr, out[0])
			curr = out
		}
	}

	c.outputs = curr
	return nil
}

func (c *Conveyer) Stop() {
	if len(c.inputs) == 0 {
		return
	}
	for _, ch := range c.inputs {
		close(ch)
	}
}
