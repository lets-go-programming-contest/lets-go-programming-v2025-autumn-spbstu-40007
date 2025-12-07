package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func safeClose(channel chan string) {
	defer func() { recover() }()
	close(channel)
}

func PrefixDecoratorFunc(
	ctx context.Context,
	inputChannel chan string,
	outputChannel chan string,
) error {
	defer safeClose(outputChannel)

	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case receivedValue, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(receivedValue, "no decorator") {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(receivedValue, prefix) {
				receivedValueWithPrefix := prefix + receivedValue
				select {
				case outputChannel <- receivedValueWithPrefix:
				case <-ctx.Done():
					return nil
				}
				continue
			}

			select {
			case outputChannel <- receivedValue:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputChannels []chan string,
	outputChannel chan string,
) error {
	defer safeClose(outputChannel)

	for {
		allInputChannelsClosed := true

		for _, inputChannel := range inputChannels {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			select {
			case receivedValue, ok := <-inputChannel:
				if !ok {
					continue
				}
				allInputChannelsClosed = false

				if strings.Contains(receivedValue, "no multiplexer") {
					continue
				}

				select {
				case outputChannel <- receivedValue:
				case <-ctx.Done():
					return nil
				}
			default:
			}
		}

		if allInputChannelsClosed {
			return nil
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	inputChannel chan string,
	outputChannels []chan string,
) error {
	if len(outputChannels) == 0 {
		return nil
	}

	defer func() {
		seenChannels := make(map[chan string]struct{})
		for _, channel := range outputChannels {
			if _, alreadySeen := seenChannels[channel]; alreadySeen {
				continue
			}
			seenChannels[channel] = struct{}{}
			safeClose(channel)
		}
	}()

	outputIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case receivedValue, ok := <-inputChannel:
			if !ok {
				return nil
			}

			targetOutputChannel := outputChannels[outputIndex%len(outputChannels)]
			outputIndex++

			select {
			case targetOutputChannel <- receivedValue:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
