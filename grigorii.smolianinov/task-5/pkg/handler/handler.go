package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	prefix := "decorated: "
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("can’t be decorated: data contains 'no decorator'")
			}

			var decoratedData string
			if strings.HasPrefix(data, prefix) {
				decoratedData = data
			} else {
				decoratedData = prefix + data
			}

			select {
			case output <- decoratedData:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			outputCh := outputs[index%len(outputs)]

			select {
			case outputCh <- data:
				index++
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	openInputChans := len(inputs)

	var mu sync.Mutex

	for _, input := range inputs {
		wg.Add(1)
		go func(input chan string) {
			defer wg.Done()
			defer func() {
				mu.Lock()
				openInputChans--
				mu.Unlock()
			}()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-input:
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
