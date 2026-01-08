package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
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

			for i := 0; i < len(outputs); i++ {
				index := (counter + i) % len(outputs)
				select {
				case <-ctx.Done():
					return nil
				case outputs[index] <- data:
					counter++
					break
				default:
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

	done := make(chan struct{})
	defer close(done)

	messages := make(chan string, len(inputs)*10)

	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(inputChan chan string) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case data, ok := <-inputChan:
					if !ok {
						return
					}
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case <-done:
						return
					case <-ctx.Done():
						return
					case messages <- data:
					}
				}
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(messages)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-messages:
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

func SimpleMultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	for _, in := range inputs {
		go func(inputChan chan string) {
			for data := range inputChan {
				if strings.Contains(data, "no multiplexer") {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- data:
				}
			}
		}(in)
	}

	<-ctx.Done()
	return nil
}
