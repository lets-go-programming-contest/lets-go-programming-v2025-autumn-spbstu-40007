package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound error = errors.New("error: chan not found")

type Conveyer interface {
	RegisterDecorator(
		fnHandler func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fnHandler func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fnHandler func(ctx context.Context, input chan string, outputs []chan string) error,
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

func New(size int) *ConveyerImpl {
	return &ConveyerImpl{
		channels:   make(map[string]chan string),
		bufferSize: size,
		handlers:   make([]func(ctx context.Context) error, 0),
	}
}

func (conveyer *ConveyerImpl) getOrCreateChan(index string) chan string {
	channel, ok := conveyer.channels[index]
	if ok {
		return channel
	}

	channel = make(chan string, conveyer.bufferSize)

	conveyer.channels[index] = channel

	return channel
}

func (conveyer *ConveyerImpl) RegisterDecorator(
	fnHandler func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputCh := conveyer.getOrCreateChan(input)
	outputCh := conveyer.getOrCreateChan(output)

	decoratorHandler := func(ctx context.Context) error {
		return fnHandler(ctx, inputCh, outputCh)
	}

	conveyer.handlers = append(conveyer.handlers, decoratorHandler)
}

func (conveyer *ConveyerImpl) RegisterMultiplexer(
	fnHandler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputsCh := make([]chan string, len(inputs))

	for i, c := range inputs {
		inputsCh[i] = conveyer.getOrCreateChan(c)
	}

	outputCh := conveyer.getOrCreateChan(output)

	multiplexerHandler := func(ctx context.Context) error {
		return fnHandler(ctx, inputsCh, outputCh)
	}

	conveyer.handlers = append(conveyer.handlers, multiplexerHandler)
}

func (conveyer *ConveyerImpl) RegisterSeparator(
	fnHandler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	outputsCh := make([]chan string, len(outputs))

	for i, c := range outputs {
		outputsCh[i] = conveyer.getOrCreateChan(c)
	}

	inputCh := conveyer.getOrCreateChan(input)

	separatorHandler := func(ctx context.Context) error {
		return fnHandler(ctx, inputCh, outputsCh)
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
	channel, exists := conveyer.channels[id]
	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (conveyer *ConveyerImpl) closeAllChannels() {
	for _, channel := range conveyer.channels {
		close(channel)
	}
}

func (conveyer *ConveyerImpl) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defer conveyer.closeAllChannels()

	var wGroup sync.WaitGroup

	errChan := make(chan error, 1)

	for _, handler := range conveyer.handlers {
		wGroup.Add(1)

		tempHandler := func(h func(context.Context) error) {
			defer wGroup.Done()

			err := h(ctx)
			if err != nil {
				select {
				case errChan <- err:
					cancel()
				default:
				}
			}
		}
		go tempHandler(handler)
	}

	select {
	case err := <-errChan:
		wGroup.Wait()

		return err
	case <-ctx.Done():
		wGroup.Wait()

		return nil
	}
}
