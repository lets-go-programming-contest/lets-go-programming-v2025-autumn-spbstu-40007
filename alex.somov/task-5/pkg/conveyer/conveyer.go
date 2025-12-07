package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type decoratorDescriptor struct {
	handlerFunction func(ctx context.Context, inputChannel chan string, outputChannel chan string) error
	inputChannelID  string
	outputChannelID string
}

type multiplexerDescriptor struct {
	handlerFunction func(ctx context.Context, inputChannels []chan string, outputChannel chan string) error
	inputChannelIDs []string
	outputChannelID string
}

type separatorDescriptor struct {
	handlerFunction  func(ctx context.Context, inputChannel chan string, outputChannels []chan string) error
	inputChannelID   string
	outputChannelIDs []string
}

type Conveyer struct {
	channelMap map[string]chan string

	decoratorDescriptors   []decoratorDescriptor
	multiplexerDescriptors []multiplexerDescriptor
	separatorDescriptors   []separatorDescriptor

	bufferSize int
	mutex      sync.Mutex
	waitGroup  sync.WaitGroup
}

func New(bufferSize int) *Conveyer {
	return &Conveyer{
		channelMap:             make(map[string]chan string),
		decoratorDescriptors:   make([]decoratorDescriptor, 0),
		multiplexerDescriptors: make([]multiplexerDescriptor, 0),
		separatorDescriptors:   make([]separatorDescriptor, 0),
		bufferSize:             bufferSize,
		mutex:                  sync.Mutex{},
		waitGroup:              sync.WaitGroup{},
	}
}

func (c *Conveyer) ensureChannelExists(channelName string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.channelMap[channelName]; !exists {
		c.channelMap[channelName] = make(chan string, c.bufferSize)
	}
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunction func(ctx context.Context, inputChannel chan string, outputChannel chan string) error,
	inputChannelID string,
	outputChannelID string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.ensureChannelExists(inputChannelID)
	c.ensureChannelExists(outputChannelID)

	c.decoratorDescriptors = append(c.decoratorDescriptors, decoratorDescriptor{
		handlerFunction: decoratorFunction,
		inputChannelID:  inputChannelID,
		outputChannelID: outputChannelID,
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunction func(ctx context.Context, inputChannels []chan string, outputChannel chan string) error,
	inputChannelIDs []string,
	outputChannelID string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, inputChannelID := range inputChannelIDs {
		c.ensureChannelExists(inputChannelID)
	}

	c.ensureChannelExists(outputChannelID)

	c.multiplexerDescriptors = append(c.multiplexerDescriptors, multiplexerDescriptor{
		handlerFunction: multiplexerFunction,
		inputChannelIDs: append([]string(nil), inputChannelIDs...),
		outputChannelID: outputChannelID,
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunction func(ctx context.Context, inputChannel chan string, outputChannels []chan string) error,
	inputChannelID string,
	outputChannelIDs []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.ensureChannelExists(inputChannelID)

	for _, outputChannelID := range outputChannelIDs {
		c.ensureChannelExists(outputChannelID)
	}

	c.separatorDescriptors = append(c.separatorDescriptors, separatorDescriptor{
		handlerFunction:  separatorFunction,
		inputChannelID:   inputChannelID,
		outputChannelIDs: append([]string(nil), outputChannelIDs...),
	})
}

func (c *Conveyer) Send(inputChannelName string, data string) error {
	c.mutex.Lock()
	inputChannel, exists := c.channelMap[inputChannelName]
	c.mutex.Unlock()

	if !exists {
		return ErrChanNotFound
	}

	inputChannel <- data
	return nil
}

func (c *Conveyer) Recv(outputChannelName string) (string, error) {
	c.mutex.Lock()
	outputChannel, exists := c.channelMap[outputChannelName]
	c.mutex.Unlock()

	if !exists {
		return "", ErrChanNotFound
	}

	receivedValue, channelClosed := <-outputChannel

	if channelClosed {
		return "undefined", nil
	}

	return receivedValue, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errorChannel := make(chan error, 1)

	c.mutex.Lock()
	decoratorDescriptors := make([]decoratorDescriptor, len(c.decoratorDescriptors))
	copy(decoratorDescriptors, c.decoratorDescriptors)

	multiplexerDescriptors := make([]multiplexerDescriptor, len(c.multiplexerDescriptors))
	copy(multiplexerDescriptors, c.multiplexerDescriptors)

	separatorDescriptors := make([]separatorDescriptor, len(c.separatorDescriptors))
	copy(separatorDescriptors, c.separatorDescriptors)
	c.mutex.Unlock()

	for _, descriptor := range decoratorDescriptors {
		c.waitGroup.Add(1)
		go func(descriptorCopy decoratorDescriptor) {
			defer c.waitGroup.Done()
			inputChannel := c.channelMap[descriptorCopy.inputChannelID]
			outputChannel := c.channelMap[descriptorCopy.outputChannelID]

			if err := descriptorCopy.handlerFunction(ctx, inputChannel, outputChannel); err != nil {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}(descriptor)
	}

	for _, descriptor := range multiplexerDescriptors {
		c.waitGroup.Add(1)
		go func(descriptorCopy multiplexerDescriptor) {
			defer c.waitGroup.Done()
			var inputChannels []chan string
			for _, inputChannelID := range descriptorCopy.inputChannelIDs {
				inputChannels = append(inputChannels, c.channelMap[inputChannelID])
			}
			outputChannel := c.channelMap[descriptorCopy.outputChannelID]

			if err := descriptorCopy.handlerFunction(ctx, inputChannels, outputChannel); err != nil {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}(descriptor)
	}

	for _, descriptor := range separatorDescriptors {
		c.waitGroup.Add(1)
		go func(descriptorCopy separatorDescriptor) {
			defer c.waitGroup.Done()
			inputChannel := c.channelMap[descriptorCopy.inputChannelID]
			var outputChannels []chan string
			for _, outputChannelID := range descriptorCopy.outputChannelIDs {
				outputChannels = append(outputChannels, c.channelMap[outputChannelID])
			}

			if err := descriptorCopy.handlerFunction(ctx, inputChannel, outputChannels); err != nil {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}(descriptor)
	}

	select {
	case receivedError := <-errorChannel:
		cancel()
		c.waitGroup.Wait()
		return receivedError
	case <-ctx.Done():
		c.waitGroup.Wait()
		return nil
	}
}
