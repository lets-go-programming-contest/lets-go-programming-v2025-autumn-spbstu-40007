package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("%w", ErrCantBeDecorated)
			}
			result := data
			if !strings.HasPrefix(data, prefix) {
				result = prefix + data
			}
			select {
			case output <- result:
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

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			target := outputs[index%len(outputs)]
			index++
			select {
			case target <- data:
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

	waitGroup := sync.WaitGroup{}
	dataChannel := make(chan string, len(inputs))

	for _, inputChannel := range inputs {
		inputChan := inputChannel
		waitGroup.Add(1)
		go func() {
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
					case dataChannel <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	go func() {
		waitGroup.Wait()
		close(dataChannel)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-dataChannel:
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
