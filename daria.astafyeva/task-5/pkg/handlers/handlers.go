package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCannotDecorate   = errors.New("can't be decorated")
	ErrNoOutputChannels = errors.New("no output channels")
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(data, "no decorator") {
				return ErrCannotDecorate
			}
			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}
			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, input := range inputs {
		go func(ch <-chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return
					}
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(input)
	}

	wg.Wait()
	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrNoOutputChannels
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			select {
			case outputs[index%len(outputs)] <- data:
				index++
			case <-ctx.Done():
				return nil
			}
		}
	}
}
