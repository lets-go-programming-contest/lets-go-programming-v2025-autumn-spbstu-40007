package conveyer

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type decoratorEntry struct {
	fn     func(ctx context.Context, input chan string, output chan string) error
	input  string
	output string
}

type separatorEntry struct {
	fn      func(ctx context.Context, input chan string, outputs []chan string) error
	input   string
	outputs []string
}

type multiplexerEntry struct {
	fn     func(ctx context.Context, inputs []chan string, output chan string) error
	inputs []string
	output string
}

type pipeline struct {
	size int

	chans map[string]chan string

	decorators   []decoratorEntry
	separators   []separatorEntry
	multiplexers []multiplexerEntry

	mu sync.RWMutex
}

func New(size int) *pipeline {
	return &pipeline{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (p *pipeline) getOrCreate(id string) chan string {
	ch, ok := p.chans[id]
	if !ok {
		ch = make(chan string, p.size)
		p.chans[id] = ch
	}
	return ch
}

func (p *pipeline) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.getOrCreate(input)
	p.getOrCreate(output)

	p.decorators = append(p.decorators, decoratorEntry{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (p *pipeline) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, id := range inputs {
		p.getOrCreate(id)
	}
	p.getOrCreate(output)

	p.multiplexers = append(p.multiplexers, multiplexerEntry{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (p *pipeline) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.getOrCreate(input)
	for _, id := range outputs {
		p.getOrCreate(id)
	}

	p.separators = append(p.separators, separatorEntry{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (p *pipeline) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, d := range p.decorators {
		d := d
		g.Go(func() error {
			return d.fn(ctx, p.chans[d.input], p.chans[d.output])
		})
	}

	for _, s := range p.separators {
		s := s
		g.Go(func() error {
			outs := make([]chan string, len(s.outputs))
			for i, id := range s.outputs {
				outs[i] = p.chans[id]
			}
			return s.fn(ctx, p.chans[s.input], outs)
		})
	}

	for _, m := range p.multiplexers {
		m := m
		g.Go(func() error {
			ins := make([]chan string, len(m.inputs))
			for i, id := range m.inputs {
				ins[i] = p.chans[id]
			}
			return m.fn(ctx, ins, p.chans[m.output])
		})
	}

	return g.Wait()
}

func (p *pipeline) Send(id string, data string) error {
	p.mu.RLock()
	ch, ok := p.chans[id]
	p.mu.RUnlock()

	if !ok {
		return errors.New("chan not found")
	}

	select {
	case ch <- data:
		return nil
	default:
		select {
		case ch <- data:
			return nil
		case <-time.After(50 * time.Millisecond):
			return errors.New("send blocked")
		}
	}
}

func (p *pipeline) Recv(id string) (string, error) {
	p.mu.RLock()
	ch, ok := p.chans[id]
	p.mu.RUnlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return val, nil
}
