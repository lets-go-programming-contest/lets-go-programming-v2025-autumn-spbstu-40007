package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound error = errors.New("error: chan not found")

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type ConveyerImpl struct {
	channels   map[string]chan string
	bufferSize int
	handlers   []func(ctx context.Context) error
}

func New(size int) Conveyer {
	return &ConveyerImpl{
		channels:   make(map[string]chan string),
		bufferSize: size,
		handlers:   make([]func(ctx context.Context) error, 0),
	}
}

func (conveyer *ConveyerImpl) getOrCreateChan(id string) chan string {
	channel, ok := conveyer.channels[id]
	if ok {
		return channel
	}

	channel = make(chan string, conveyer.bufferSize)

	conveyer.channels[id] = channel

	return channel
}

func (conveyer *ConveyerImpl) closeAllChannels() {
	for _, channel := range conveyer.channels {
		close(channel)
	}

	conveyer.channels = make(map[string]chan string)
}

func (conveyer *ConveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputCh := conveyer.getOrCreateChan(input)
	outputCh := conveyer.getOrCreateChan(output)

	decoratorHandler := func(ctx context.Context) error {
		return fn(ctx, inputCh, outputCh)
	}

	conveyer.handlers = append(conveyer.handlers, decoratorHandler)
}

func (conveyer *ConveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputsCh := make([]chan string, len(inputs))

	for i, c := range inputs {
		inputsCh[i] = conveyer.getOrCreateChan(c)
	}

	outputCh := conveyer.getOrCreateChan(output)

	multiplexerHandler := func(ctx context.Context) error {
		return fn(ctx, inputsCh, outputCh)
	}

	conveyer.handlers = append(conveyer.handlers, multiplexerHandler)
}

func (conveyer *ConveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	outputsCh := make([]chan string, len(outputs))

	for i, c := range outputs {
		outputsCh[i] = conveyer.getOrCreateChan(c)
	}

	inputCh := conveyer.getOrCreateChan(input)

	separatorHandler := func(ctx context.Context) error {
		return fn(ctx, inputCh, outputsCh)
	}

	conveyer.handlers = append(conveyer.handlers, separatorHandler)
}

func (conveyer *ConveyerImpl) Send(id string, data string) error {
	channel, ok := conveyer.channels[id]
	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (conveyer *ConveyerImpl) Recv(id string) (string, error) {
	channel, ok := conveyer.channels[id]
	if !ok {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (conveyer *ConveyerImpl) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	done := make(chan struct{})

	for _, handler := range conveyer.handlers {
		h := handler
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h(ctx); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:
		cancel()
		wg.Wait()
		conveyer.closeAllChannels()
		return err
	case <-done:
		conveyer.closeAllChannels()
		return nil
	case <-ctx.Done():
		cancel()
		wg.Wait()
		conveyer.closeAllChannels()
		return ctx.Err()
	}
}
