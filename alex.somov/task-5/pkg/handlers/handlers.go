package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func safeClose(ch chan string) {
	defer func() { _ = recover() }()
	close(ch)
}

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer safeClose(output)

	const prefix = "decorated: "

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

			if !strings.HasPrefix(value, prefix) {
				value = prefix + value
			}

			select {
			case output <- value:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func copyInput(ctx context.Context, input <-chan string, out chan<- string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-input:
			if !ok {
				return
			}
			select {
			case out <- v:
			case <-ctx.Done():
				return
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer safeClose(output)

	relay := make(chan string)

	var WaitGroup sync.WaitGroup

	for _, in := range inputs {
		WaitGroup.Add(1)

		go copyInput(ctx, in, relay, &WaitGroup)
	}

	go func() {
		WaitGroup.Wait()
		close(relay)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-relay:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no multiplexer") {
				continue
			}

			select {
			case output <- value:
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

	defer func() {
		seen := make(map[chan string]struct{})
		for _, channel := range outputs {
			if _, ok := seen[channel]; ok {
				continue
			}

			seen[channel] = struct{}{}

			safeClose(channel)
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

			out := outputs[idx%len(outputs)]
			idx++

			select {
			case out <- value:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
