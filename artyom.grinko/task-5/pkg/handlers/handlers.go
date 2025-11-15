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

	for {
		select {
		case <-context.Done():
			return nil

		default:
			for _, input := range inputs {
				select {
				case <-context.Done():
					return nil

				case x, ok := <-input:
					if !ok {
						return nil
					}

					if !strings.Contains(x, "no multiplexer") {
						select {
						case <-context.Done():
							return nil

						case output <- x:
						}
					}
				}
			}
		}
	}
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
