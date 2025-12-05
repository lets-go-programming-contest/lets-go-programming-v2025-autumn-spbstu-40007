package handlers

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

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

			var result string
			if strings.HasPrefix(data, prefix) {
				result = data
			} else {
				result = prefix + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- result:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := atomic.AddUint64(&counter, 1) - 1
			outputIdx := int(idx % uint64(len(outputs)))

			select {
			case <-ctx.Done():
				return nil
			case outputs[outputIdx] <- data:
			}
		}
	}
}

func multiplexerWorker(ctx context.Context, input <-chan string, output chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	merged := make(chan string, len(inputs)*2)

	for _, inputChan := range inputs {
		go func(in <-chan string) {
			_ = multiplexerWorker(ctx, in, merged)
		}(inputChan)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-merged:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}
