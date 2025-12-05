package handlers

import (
	"context"
	"strings"
	"sync"
)

type DecoratingHandler func(ctx context.Context, in chan string, out chan string) error
type MultiplexingHandler func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatingHandler func(ctx context.Context, input chan string, outputs []chan string) error

func PrefixDecoratorFunc(ctx context.Context, in chan string, out chan string) error {
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}
			if strings.Contains(data, "no decorator") {
				continue
			}
			select {
			case <-ctx.Done():
				return nil
			case out <- "decorated: " + data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, in := range inputs {
		ch := in
		go func() {
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
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

func SeparatorFunc(ctx context.Context, in chan string, outs []chan string) error {
	defer func() {
		for _, ch := range outs {
			close(ch)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}

			idx := 0
			if strings.Contains(data, "separator B") {
				idx = 1
			}

			select {
			case <-ctx.Done():
				return nil
			case outs[idx] <- data:
			}
		}
	}
}
