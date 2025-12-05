package conveyer

import (
	"context"
	"errors"

	"github.com/ksuah/task-5/pkg/handlers"
)

type step struct {
	dec handlers.DecoratingHandler
	mux handlers.MultiplexingHandler
	isMux bool
}

type Conveyer struct {
	steps   []step
	ctx     context.Context
	cancel  context.CancelFunc

	input  chan string
	output chan string
	queueSize int
}

func New(queueSize int) *Conveyer {
	return &Conveyer{
		queueSize: queueSize,
	}
}

func (c *Conveyer) RegisterDecorator(name string, h handlers.DecoratingHandler) {
	c.steps = append(c.steps, step{
		dec: h,
	})
}

func (c *Conveyer) RegisterMultiplexer(name string, h handlers.MultiplexingHandler) {
	c.steps = append(c.steps, step{
		mux: h,
		isMux: true,
	})
}

func (c *Conveyer) Send(data string) error {
	if c.input == nil {
		return errors.New("conveyer not started")
	}
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case c.input <- data:
		return nil
	}
}

func (c *Conveyer) Recv() (string, error) {
	if c.output == nil {
		return "", errors.New("conveyer not started")
	}
	select {
	case <-c.ctx.Done():
		return "", nil
	case v := <-c.output:
		return v, nil
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	currInputs := []chan string{make(chan string, c.queueSize)}
	c.input = currInputs[0]

	for _, s := range c.steps {
		if !s.isMux {
			out := make(chan string, c.queueSize)
			go s.dec(c.ctx, currInputs[0], out)
			currInputs = []chan string{out}
		} else {
			out := make(chan string, c.queueSize)
			go s.mux(c.ctx, currInputs, out)
			currInputs = []chan string{out}
		}
	}

	c.output = currInputs[0]

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
