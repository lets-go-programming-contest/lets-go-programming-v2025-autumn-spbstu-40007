package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(prefix string) func(
	ctx context.Context, input chan string, output chan string,
) error {
	return func(ctx context.Context, input chan string, output chan string) error {
		defer close(output)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case value, ok := <-input:
				if !ok {
					return nil
				}

				if value == "" {
					return ErrCannotDecorate
				}

				output <- prefix + value
			}
		}
	}
}

func SeparatorFunc(sep string) func(
	ctx context.Context, input chan string, outputs []chan string,
) error {
	return func(ctx context.Context, input chan string, outputs []chan string) error {
		defer func() {
			for _, out := range outputs {
				close(out)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case value, ok := <-input:
				if !ok {
					return nil
				}

				parts := strings.Split(value, sep)

				for idx, part := range parts {
					if idx < len(outputs) {
						outputs[idx] <- part
					}
				}
			}
		}
	}
}

func MultiplexerFunc() func(
	ctx context.Context, inputs []chan string, output chan string,
) error {
	return func(ctx context.Context, inputs []chan string, output chan string) error {
		defer close(output)

		done := make(chan struct{})
		active := len(inputs)

		for _, inputCh := range inputs {
			go func(input chan string) {
				for {
					select {
					case <-ctx.Done():
						done <- struct{}{}

						return
					case value, ok := <-input:
						if !ok {
							done <- struct{}{}

							return
						}

						output <- value
					}
				}
			}(inputCh)
		}

		for active > 0 {
			<-done

			active--
		}

		return nil
	}
}
