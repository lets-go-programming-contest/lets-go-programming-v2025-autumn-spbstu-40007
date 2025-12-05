package conveyer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

const emptyPayload = "undefined"

var ErrStreamNotFound = errors.New("chan not found")

type Pipeline struct {
	bufferSize int
	streams    map[string]chan string
	workers    []func(context.Context) error
}

func New(bufferSize int) *Pipeline {
	return &Pipeline{
		bufferSize: bufferSize,
		streams:    make(map[string]chan string),
		workers:    make([]func(context.Context) error, 0),
	}
}

func (p *Pipeline) makeStream(name string) {
	if _, ok := p.streams[name]; !ok {
		p.streams[name] = make(chan string, p.bufferSize)
	}
}

func (p *Pipeline) makeStreams(names ...string) {
	for _, n := range names {
		p.makeStream(n)
	}
}

func (p *Pipeline) RegisterDecorator(
	handler func(ctx context.Context, source chan string, sink chan string) error,
	source string,
	sink string,
) {
	p.makeStreams(source)
	p.makeStreams(sink)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, p.streams[source], p.streams[sink])
	})
}

func (p *Pipeline) RegisterMultiplexer(
	handler func(ctx context.Context, sources []chan string, sink chan string) error,
	sourceNames []string,
	sink string,
) {
	p.makeStreams(sourceNames...)
	p.makeStreams(sink)

	p.workers = append(p.workers, func(ctx context.Context) error {
		sourceStreams := make([]chan string, 0, len(sourceNames))
		for _, name := range sourceNames {
			sourceStreams = append(sourceStreams, p.streams[name])
		}

		return handler(ctx, sourceStreams, p.streams[sink])
	})
}

func (p *Pipeline) RegisterSeparator(
	handler func(ctx context.Context, source chan string, sinks []chan string) error,
	source string,
	sinkNames []string,
) {
	p.makeStreams(source)
	p.makeStreams(sinkNames...)

	p.workers = append(p.workers, func(ctx context.Context) error {
		sinkStreams := make([]chan string, 0, len(sinkNames))
		for _, name := range sinkNames {
			sinkStreams = append(sinkStreams, p.streams[name])
		}

		return handler(ctx, p.streams[source], sinkStreams)
	})
}

func (p *Pipeline) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, handlerFunc := range p.workers {
		workerFunc := handlerFunc

		group.Go(func() error {
			return workerFunc(groupCtx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("pipeline workers: %w", err)
	}

	return nil
}

func (p *Pipeline) Send(source string, payload string) error {
	stream, ok := p.streams[source]
	if !ok {
		return ErrStreamNotFound
	}

	stream <- payload

	return nil
}

func (p *Pipeline) Recv(sink string) (string, error) {
	stream, ok := p.streams[sink]
	if !ok {
		return "", ErrStreamNotFound
	}

	receivedPayload, open := <-stream

	if !open {
		return emptyPayload, nil
	}

	return receivedPayload, nil
}
