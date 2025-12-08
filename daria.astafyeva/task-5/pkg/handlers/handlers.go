package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCannotDecorate   = errors.New("can't be decorated")
	ErrNoOutputChannels = errors.New("no output channels")
)

func PrefixDecoratorFunc(ctx context.Context, inputChannel, outputChannel chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case outputChannel <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	var waitGroup sync.WaitGroup

	for _, channel := range inputChannels {
		waitGroup.Add(1)

		go func(inputChan <-chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inputChan:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case outputChannel <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(channel)
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, inputChannel chan string, outputChannels []chan string) error {
	if len(outputChannels) == 0 {
		return ErrNoOutputChannels
	}

	position := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			target := outputChannels[position%len(outputChannels)]
			select {
			case target <- data:
				position++
			case <-ctx.Done():
				return nil
			}
		}
	}
}
