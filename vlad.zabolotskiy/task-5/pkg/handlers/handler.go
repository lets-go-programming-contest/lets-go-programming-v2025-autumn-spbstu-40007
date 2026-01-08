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
			counter++

			select {
			case <-ctx.Done():
				return nil
			case outputs[index] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	for {
		anyActive := false

		for _, in := range inputs {
			select {
			case <-ctx.Done():
				return nil
			case data, ok := <-in:
				if ok {
					anyActive = true
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case <-ctx.Done():
						return nil
					case output <- data:
					}
				}
			default:
				anyActive = true
			}
		}

		if !anyActive {
			return nil
		}
	}
}
