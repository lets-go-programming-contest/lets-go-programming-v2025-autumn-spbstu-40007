package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	const prefixText = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("%w", ErrCannotBeDecorated)
			}

			resultData := data
			if !strings.HasPrefix(data, prefixText) {
				resultData = prefixText + data
			}

			select {
			case output <- resultData:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	outputIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			targetOutput := outputs[outputIndex%len(outputs)]
			outputIndex++

			select {
			case targetOutput <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	inputWaitGroup := sync.WaitGroup{}
	internalDataChannel := make(chan string, len(inputs))

	inputWaitGroup.Add(len(inputs))
	for _, inputChannel := range inputs {
		specificInput := inputChannel
		inputWaitGroup.Add(1)
		go func(inputChan chan string) {
			defer inputWaitGroup.Done()
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
					case internalDataChannel <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(specificInput)
	}

	go func() {
		inputWaitGroup.Wait()
		close(internalDataChannel)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-internalDataChannel:
			if !ok {
				return nil
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
