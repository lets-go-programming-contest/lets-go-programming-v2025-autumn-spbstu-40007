package handlers

import (
	"context"
	"strings"
)

func PrefixDecoratorFunc(prefix string) func(ctx context.Context, input chan string, output chan string, errCh chan error) {
	return func(ctx context.Context, input chan string, output chan string, errCh chan error) {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-input:
				if !ok {
					return
				}
				output <- prefix + v
			}
		}
	}
}

func SeparatorFunc(sep string) func(ctx context.Context, input chan string, outputs []chan string, errCh chan error) {
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
			case v, ok := <-input:
				if !ok {
					return
				}
				parts := strings.Split(v, sep)
				for i, part := range parts {
					if i < len(outputs) {
						outputs[i] <- part
					}
				}
			}
		}
	}
}

func MultiplexerFunc() func(ctx context.Context, inputs []chan string, output chan string, errCh chan error) {
	return func(ctx context.Context, inputs []chan string, output chan string, errCh chan error) {
		defer close(output)

		done := make(chan struct{})
		active := len(inputs)

		for _, input := range inputs {
			in := input
			go func() {
				for {
					select {
					case <-ctx.Done():
						done <- struct{}{}
						return
					case v, ok := <-in:
						if !ok {
							done <- struct{}{}
							return
						}
						output <- v
					}
				}
			}()
		}

		for active > 0 {
			<-done
			active--
		}
	}
}
