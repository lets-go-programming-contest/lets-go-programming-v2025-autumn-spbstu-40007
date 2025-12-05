package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(item, "no decorator") {
				return errors.New("can't be decorated")
			}

			prefix := "decorated: "
			if !strings.HasPrefix(item, prefix) {
				item = prefix + item
			}

			select {
			case output <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	merge := func(in chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-in:
				if !ok {
					return
				}

				if strings.Contains(item, "no multiplexer") {
					continue
				}

				select {
				case output <- item:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, in := range inputs {
		wg.Add(1)
		go merge(in)
	}

	wg.Wait()
	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	var i int
	count := len(outputs)
	if count == 0 {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}

			targetIndex := i % count
			i++

			select {
			case outputs[targetIndex] <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
