package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated error = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
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

			stringDecorated := data
			if !strings.Contains(data, "decorated: ") {
				stringDecorated = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- stringDecorated:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
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
	var wGroup sync.WaitGroup

	for _, input := range inputs {
		wGroup.Add(1)

		multiplexerProcess := func(inChan chan string) {
			defer wGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inChan:
					if !ok {
						return
					}

					if !strings.Contains(data, "no multiplexer") {
						select {
						case output <- data:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}
		go multiplexerProcess(input)
	}

	wGroup.Wait()

	return nil
}
