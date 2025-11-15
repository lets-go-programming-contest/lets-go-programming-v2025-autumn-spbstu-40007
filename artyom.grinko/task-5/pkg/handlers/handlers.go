package handlers

import (
	"context"
	"fmt"
	"strings"
)

func PrefixDecoratorFunc(
	context context.Context,
	input chan string,
	output chan string,
) error {
	defer close(output)

	for {
		select {
		case <-context.Done():
			return nil

		case x, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(x, "no decorator") {
				return fmt.Errorf("handlers.PrefixDecoratorFunc: can't be decorated")
			}

			if strings.HasPrefix(x, "decorated: ") {
				output <- x

				continue
			}

			select {
			case <-context.Done():
				return nil

			case output <- "decorated: " + x:
			}
		}
	}
}

func MultiplexerFunc(
	context context.Context,
	inputs []chan string,
	output chan string,
) error {
	defer close(output)

	activeInputs := make([]chan string, len(inputs))
	copy(activeInputs, inputs)

	for len(activeInputs) > 0 {
		select {
		case <-context.Done():
			return nil

		default:
			for i := 0; i < len(activeInputs); i++ {
				select {
				case <-context.Done():
					return nil

				case x, ok := <-activeInputs[i]:
					if !ok {
						activeInputs = append(activeInputs[:i], activeInputs[i+1:]...)
						i--

						continue
					}

					if strings.Contains(x, "no multiplexer") {
						continue
					}

					select {
					case <-context.Done():
						return nil

					case output <- x:
					}
				default:
				}
			}
		}
	}

	return nil
}

func SeparatorFunc(
	context context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer (func() {
		for _, output := range outputs {
			close(output)
		}
	})()

	for i := 0; ; {
		select {
		case <-context.Done():
			return nil

		case x, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-context.Done():
				return nil

			case outputs[i%len(outputs)] <- x:
				i++
			}
		}
	}
}
