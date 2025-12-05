package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

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
				return ErrCannotBeDecorated
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
	var waitGroup sync.WaitGroup

	merge := func(inputChannel chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-inputChannel:
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

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		go merge(inputChannel)
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	var index int

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

			targetIndex := index % count
			index++

			select {
			case outputs[targetIndex] <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
