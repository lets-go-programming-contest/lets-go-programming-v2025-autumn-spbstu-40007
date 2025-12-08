package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return ErrCannotDecorate
			}

			decoratedValue := "decorated: " + value

			select {
			case <-ctx.Done():
				return nil
			case output <- decoratedValue:
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

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}

			if len(outputs) > 0 {
				select {
				case <-ctx.Done():
					return nil
				case outputs[idx%len(outputs)] <- value:
					idx++
				}
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	for _, inputCh := range inputs {
		localInputCh := inputCh
		wg.Add(1)

		go func(ch chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(value, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- value:
					}
				}
			}
		}(localInputCh)
	}

	wg.Wait()

	return nil
}
