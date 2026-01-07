package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			index := counter % len(outputs)

			select {
			case <-ctx.Done():
				return nil
			case outputs[index] <- data:
				counter++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	activeInputs := make([]chan string, len(inputs))
	copy(activeInputs, inputs)

	for len(activeInputs) > 0 {
		for index := 0; index < len(activeInputs); index++ {
			select {
			case <-ctx.Done():
				return nil

			case data, ok := <-activeInputs[index]:
				if !ok {
					activeInputs = append(activeInputs[:index], activeInputs[index+1:]...)
					index--

					continue
				}

				if strings.Contains(data, "no multiplexer") {
					continue
				}

				select {
				case <-ctx.Done():
					return nil
				case output <- data:
				}
			default:
			}
		}
	}

	return nil
}
