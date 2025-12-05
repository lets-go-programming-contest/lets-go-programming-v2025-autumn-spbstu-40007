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
	tp       StepType
	name     string
	input    string
	inputs   []string
	outputs  []string
	dec      handlers.DecoratingHandler
	sep      handlers.SeparatingHandler
	mux      handlers.MultiplexingHandler
}

type Conveyer struct {
	steps   []step
	ctx     context.Context
	cancel  context.CancelFunc
	inputs  map[string]chan string
	outputs map[string]chan string
}

func New(names ...string) *Conveyer {
	inputs := make(map[string]chan string)
	for _, n := range names {
		inputs[n] = make(chan string)
	}
	return &Conveyer{
		inputs:  inputs,
		outputs: make(map[string]chan string),
	}
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

func (c *Conveyer) RegisterDecorator(name, inputName, outputName string, h handlers.DecoratingHandler) {
	c.steps = append(c.steps, step{
		tp:      Decorating,
		name:    name,
		input:   inputName,
		outputs: []string{outputName},
		dec:     h,
	})
}

func (c *Conveyer) RegisterSeparator(name, inputName string, outputNames []string, h handlers.SeparatingHandler) {
	c.steps = append(c.steps, step{
		tp:      Separating,
		name:    name,
		input:   inputName,
		outputs: outputNames,
		sep:     h,
	})
}

func (c *Conveyer) RegisterMultiplexer(name string, inputNames []string, outputName string, h handlers.MultiplexingHandler) {
	c.steps = append(c.steps, step{
		tp:      Multiplexing,
		name:    name,
		inputs:  inputNames,
		outputs: []string{outputName},
		mux:     h,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	channels := make(map[string]chan string)
	for _, s := range c.steps {
		for _, out := range s.outputs {
			channels[out] = make(chan string)
		}
	}
	for _, s := range c.steps {
		switch s.tp {
		case Decorating:
			in := c.inputs[s.input]
			if in == nil {
				in = channels[s.input]
			}
			out := channels[s.outputs[0]]
			go s.dec(c.ctx, in, out)
		case Separating:
			in := c.inputs[s.input]
			if in == nil {
				in = channels[s.input]
			}
			outs := make([]chan string, len(s.outputs))
			for i, n := range s.outputs {
				outs[i] = channels[n]
			}
			go s.sep(c.ctx, in, outs)
		case Multiplexing:
			ins := make([]chan string, len(s.inputs))
			for i, n := range s.inputs {
				ch := c.inputs[n]
				if ch == nil {
					ch = channels[n]
				}
				ins[i] = ch
			}
			out := channels[s.outputs[0]]
			go s.mux(c.ctx, ins, out)
		}
	}
	for name, ch := range channels {
		c.outputs[name] = ch
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
	for _, ch := range c.outputs {
		close(ch)
	}
}
