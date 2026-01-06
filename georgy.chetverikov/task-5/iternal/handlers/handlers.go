package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
)

const prefix = "decorated: "

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}

			var result string
			if strings.HasPrefix(data, prefix) {
				result = data
			} else {
				result = prefix + data
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- result:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var counter int64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := atomic.AddInt64(&counter, 1) - 1
			targetIdx := int(idx % int64(len(outputs)))

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[targetIdx] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	errChan := make(chan error, len(inputs))
	done := make(chan struct{})

	for _, input := range inputs {
		go func(in chan string) {
			for {
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				case <-done:
					return
				case data, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						errChan <- ctx.Err()
						return
					case <-done:
						return
					case output <- data:
					}
				}
			}
		}(input)
	}

	select {
	case <-ctx.Done():
		close(done)
		return ctx.Err()
	case err := <-errChan:
		close(done)
		return err
	}
}
