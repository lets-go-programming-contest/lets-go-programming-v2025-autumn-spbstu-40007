package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type ConveyerImpl struct {
	channelSize int
	channels    map[string]chan string
	runners     []func(ctx context.Context) error
	mu          sync.RWMutex
}

func New(size int) *ConveyerImpl {
	return &ConveyerImpl{
		channelSize: size,
		channels:    make(map[string]chan string),
		runners:     make([]func(ctx context.Context) error, 0),
		mu:          sync.RWMutex{},
	}
}

func (c *ConveyerImpl) getOrCreateChannel(channelID string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[channelID]; ok {
		return ch
	}

	channel := make(chan string, c.channelSize)
	c.channels[channelID] = channel

	return channel
}

func (c *ConveyerImpl) RegisterDecorator(
	fnRunner func(ctx context.Context, input chan string, output chan string) error,
	inputID string,
	outputID string,
) {
	inputCh := c.getOrCreateChannel(inputID)
	outputCh := c.getOrCreateChannel(outputID)

	runner := func(ctx context.Context) error {
		return fnRunner(ctx, inputCh, outputCh)
	}

	c.mu.Lock()
	c.runners = append(c.runners, runner)
	c.mu.Unlock()
}

func (c *ConveyerImpl) RegisterMultiplexer(
	fnRunner func(ctx context.Context, inputs []chan string, output chan string) error,
	inputIDs []string,
	outputID string,
) {
	inputChans := make([]chan string, len(inputIDs))
	for idx, id := range inputIDs {
		inputChans[idx] = c.getOrCreateChannel(id)
	}

	outputCh := c.getOrCreateChannel(outputID)

	runner := func(ctx context.Context) error {
		return fnRunner(ctx, inputChans, outputCh)
	}

	c.mu.Lock()
	c.runners = append(c.runners, runner)
	c.mu.Unlock()
}

func (c *ConveyerImpl) RegisterSeparator(
	fnRunner func(ctx context.Context, input chan string, outputs []chan string) error,
	inputID string,
	outputIDs []string,
) {
	inputCh := c.getOrCreateChannel(inputID)

	outputChans := make([]chan string, len(outputIDs))
	for idx, id := range outputIDs {
		outputChans[idx] = c.getOrCreateChannel(id)
	}

	runner := func(ctx context.Context) error {
		return fnRunner(ctx, inputCh, outputChans)
	}

	c.mu.Lock()
	c.runners = append(c.runners, runner)
	c.mu.Unlock()
}

func (c *ConveyerImpl) runWorkers(ctx context.Context, errChan chan error) {
	c.mu.RLock()
	runners := c.runners
	c.mu.RUnlock()

	var waitGroup sync.WaitGroup

	workerCount := len(runners)

	waitGroup.Add(workerCount)

	for _, runner := range runners {
		go func(r func(ctx context.Context) error) {
			defer waitGroup.Done()

			err := r(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {
				select {
				case errChan <- err:
				default:
				}
			}
		}(runner)
	}

	waitGroup.Wait()
}

func (c *ConveyerImpl) Run(ctx context.Context) error {
	c.mu.RLock()
	numRunners := len(c.runners)
	c.mu.RUnlock()

	if numRunners == 0 {
		return nil
	}

	errChan := make(chan error, numRunners)
	done := make(chan struct{})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		c.runWorkers(ctx, errChan)
		close(done)
	}()

	var runErr error

	select {
	case <-ctx.Done():
		runErr = ctx.Err()

	case err := <-errChan:
		cancel()

		runErr = err

	case <-done:
	}

	c.mu.Lock()
	for _, channel := range c.channels {
		close(channel)
	}
	c.mu.Unlock()

	select {
	case internalErr := <-errChan:
		if runErr == nil {
			runErr = internalErr
		}
	default:
	}

	if errors.Is(runErr, context.DeadlineExceeded) || errors.Is(runErr, context.Canceled) {
		return nil
	}

	return runErr
}

func (c *ConveyerImpl) Send(inputID string, data string) error {
	c.mu.RLock()
	channel, channelExists := c.channels[inputID]
	c.mu.RUnlock()

	if !channelExists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *ConveyerImpl) Recv(outputID string) (string, error) {
	c.mu.RLock()
	channel, channelExists := c.channels[outputID]
	c.mu.RUnlock()

	if !channelExists {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return value, nil
}
