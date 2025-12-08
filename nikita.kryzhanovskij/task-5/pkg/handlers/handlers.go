package handlers

import (
	"context"
	"fmt"
	"strings"
)

func PrefixDecoratorFunc(prefix string) func(
	ctx context.Context, input chan string, output chan string, errCh chan error,
) {
	return func(ctx context.Context, input chan string, output chan string, errCh chan error) {
		defer close(output)

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-input:
				if !ok {
					return
				}

				if value == "" {
					errCh <- fmt.Errorf("can't be decorated")
					continue
				}

				output <- prefix + value
			}
		}
	}
}

func SeparatorFunc(sep string) func(
	ctx context.Context, input chan string, outputs []chan string, errCh chan error,
) {
	return func(ctx context.Context, input chan string, outputs []chan string, errCh chan error) {
		defer func() {
			for _, out := range outputs {
				close(out)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-input:
				if !ok {
					return
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
	ctx context.Context, inputs []chan string, output chan string, errCh chan error,
) {
	return func(ctx context.Context, inputs []chan string, output chan string, errCh chan error) {
		defer close(output)

		done := make(chan struct{})
		active := len(inputs)

		for idx, inputCh := range inputs {
			identifier := idx

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

						fmt.Printf("[Multiplexer-%d] Получено и отправлено: %q\n", identifier, value)
						output <- value
					}
				}
			}(inputCh)
		}

		for active > 0 {
			<-done
			active--
		}
	}
}
