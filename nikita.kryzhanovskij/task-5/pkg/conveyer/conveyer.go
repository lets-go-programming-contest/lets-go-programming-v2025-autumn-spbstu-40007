package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

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

type conveyorImpl struct {
	chans   map[string]chan string
	size    int
	workers []func(ctx context.Context) error
}

func New(size int) *conveyorImpl {
	return &conveyorImpl{
		chans:   make(map[string]chan string),
		size:    size,
		workers: make([]func(ctx context.Context) error, 0),
	}
}

func (c *conveyorImpl) getOrCreateChan(name string) chan string {
	if ch, ok := c.chans[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.chans[name] = ch

	return ch
}

func (c *conveyorImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputCh := c.getOrCreateChan(input)
	outputCh := c.getOrCreateChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputCh)
	})
}

func (c *conveyorImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputsCh := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputsCh[i] = c.getOrCreateChan(name)
	}

	outputCh := c.getOrCreateChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inputsCh, outputCh)
	})
}

func (c *conveyorImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputCh := c.getOrCreateChan(input)
	outputsCh := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputsCh[i] = c.getOrCreateChan(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputsCh)
	})
}

func (c *conveyorImpl) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	errCh := make(chan error, len(c.workers))

	// Создаём новый контекст с возможностью отмены
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Закрываем все каналы после завершения
	defer func() {
		for _, ch := range c.chans {
			close(ch)
		}
	}()

	// Запускаем все workers
	for _, worker := range c.workers {
		wg.Add(1)

		w := worker

		go func() {
			defer wg.Done()

			if err := w(ctx); err != nil {
				if !errors.Is(err, context.Canceled) {
					select {
					case errCh <- err:
						cancel()
					default:
					}
				}
			}
		}()
	}

	// Ждём либо ошибку, либо отмену контекста
	select {
	case <-ctx.Done():
		wg.Wait()
		return nil
	case err := <-errCh:
		wg.Wait()
		return err
	}
}

func (c *conveyorImpl) Send(input string, data string) error {
	ch, ok := c.chans[input]
	if !ok {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *conveyorImpl) Recv(output string) (string, error) {
	ch, ok := c.chans[output]
	if !ok {
		return "", ErrChanNotFound
	}

	val, isOpen := <-ch
	if !isOpen {
		return "undefined", nil
	}

	return val, nil
}
