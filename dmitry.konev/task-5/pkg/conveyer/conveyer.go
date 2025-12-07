package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrSendBlocked  = errors.New("send blocked")
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
	size         int
	chans        map[string]chan string
	decorators   []decoratorEntry
	separators   []separatorEntry
	multiplexers []multiplexerEntry
	mu           sync.RWMutex
}

func New(size int) *pipeline {
	return &pipeline{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (p *pipeline) getOrCreate(id string) chan string {
	ch, exists := p.chans[id]
	if !exists {
		ch = make(chan string, p.size)
		p.chans[id] = ch
	}

	return ch
}

func (p *pipeline) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input, output string) {
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

func (p *pipeline) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
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

func (p *pipeline) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
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
	group, ctx := errgroup.WithContext(ctx)

	for _, dEntry := range p.decorators {
		entry := dEntry
		group.Go(func() error {
			return wrapErr(entry.fn(ctx, p.chans[entry.input], p.chans[entry.output]))
		})
	}

	for _, sEntry := range p.separators {
		entry := sEntry
		group.Go(func() error {
			var outs []chan string
			for _, id := range entry.outputs {
				outs = append(outs, p.chans[id])
			}
			return wrapErr(entry.fn(ctx, p.chans[entry.input], outs))
		})
	}

	for _, mEntry := range p.multiplexers {
		entry := mEntry
		group.Go(func() error {
			var ins []chan string
			for _, id := range entry.inputs {
				ins = append(ins, p.chans[id])
			}
			return wrapErr(entry.fn(ctx, ins, p.chans[entry.output]))
		})
	}

	err := group.Wait()

	p.mu.Lock()
	for _, ch := range p.chans {
		close(ch)
	}
	p.mu.Unlock()

	return wrapErr(err)
}

func (p *pipeline) Send(id string, data string) error {
	p.mu.RLock()
	ch, exists := p.chans[id]
	p.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	select {
	case ch <- data:
		return nil

	default:
		select {
		case ch <- data:
			return nil

		case <-time.After(50 * time.Millisecond):
			return ErrSendBlocked
		}
	}
}

func (p *pipeline) Recv(id string) (string, error) {
	p.mu.RLock()
	ch, exists := p.chans[id]
	p.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	value, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return value, nil
}

func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w", err)
}
