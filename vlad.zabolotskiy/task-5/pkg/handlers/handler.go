package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

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
	defer func() {
		for _, out := range outputs {
			close(out)
		}
	}()

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
	defer close(output)

	for {
		allClosed := true

		for _, inputChan := range inputs {
			select {
			case <-ctx.Done():
				return nil
			case data, ok := <-inputChan:
				if !ok {
					continue
				}

				allClosed = false

				if strings.Contains(data, "no multiplexer") {
					continue
				}

				select {
				case <-ctx.Done():
					return nil
				case output <- data:
				}
			default:
				select {
				case <-ctx.Done():
					return nil
				case _, ok := <-inputChan:
					if ok {
						allClosed = false
					}
				default:
				}
			}
		}

		if allClosed {
			return nil
		}
	}
}
