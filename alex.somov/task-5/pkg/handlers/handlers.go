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
		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(v, prefix) {
				v = prefix + v
			}

			select {
			case output <- v:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func copyInput(ctx context.Context, in <-chan string, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-in:
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
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go copyInput(ctx, in, relay, &wg)
	}

	go func() {
		wg.Wait()
		close(relay)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-relay:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no multiplexer") {
				continue
			}

			select {
			case output <- v:
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
		for _, ch := range outputs {
			if _, ok := seen[ch]; ok {
				continue
			}
			seen[ch] = struct{}{}
			safeClose(ch)
		}
	}()

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-input:
			if !ok {
				return nil
			}

			out := outputs[idx%len(outputs)]
			idx++

			select {
			case out <- v:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
