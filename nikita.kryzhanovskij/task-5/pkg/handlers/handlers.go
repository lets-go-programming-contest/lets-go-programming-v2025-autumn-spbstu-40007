package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc() func(ctx context.Context, input chan string, output chan string) error {
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

				if strings.Contains(value, "no decorator") {
					return ErrCannotDecorate
				}

				select {
				case output <- "decorated: " + value:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}

func SeparatorFunc() func(ctx context.Context, input chan string, outputs []chan string) error {
	return func(ctx context.Context, input chan string, outputs []chan string) error {
		defer func() {
			for _, out := range outputs {
				close(out)
			}
		}()

		idx := 0

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case value, ok := <-input:
				if !ok {
					return nil
				}

				if idx < len(outputs) {
					select {
					case outputs[idx] <- value:
					case <-ctx.Done():
						return ctx.Err()
					}
				}

				idx = (idx + 1) % len(outputs)
			}
		}
	}
}

func MultiplexerFunc() func(ctx context.Context, inputs []chan string, output chan string) error {
	return func(ctx context.Context, inputs []chan string, output chan string) error {
		defer close(output)

		done := make(chan struct{})
		active := len(inputs)

		for _, inputCh := range inputs {
			go func(inputChannel chan string) {
				defer func() { done <- struct{}{} }()

				for {
					select {
					case <-ctx.Done():
						return
					case value, ok := <-inputChannel:
						if !ok {
							return
						}

						if strings.Contains(value, "no multiplexer") {
							continue
						}

						select {
						case output <- value:
						case <-ctx.Done():
							return
						}
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
