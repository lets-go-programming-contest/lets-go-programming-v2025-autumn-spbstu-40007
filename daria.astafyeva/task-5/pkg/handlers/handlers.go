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

const (
	skipMux = "no multiplexer"
	skipDec = "no decorator"
	prefix  = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, in chan string, out chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				close(out)
				return nil
			}
			if strings.Contains(data, skipDec) {
				return ErrCannotDecorate
			}
			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}
			select {
			case out <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, ins []chan string, out chan string) error {
	var wg sync.WaitGroup
	done := make(chan struct{})

	for _, in := range ins {
		wg.Add(1)
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
					if strings.Contains(data, skipMux) {
						continue
					}
					select {
					case out <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(done)
		close(out)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func SeparatorFunc(ctx context.Context, in chan string, outs []chan string) error {
	if len(outs) == 0 {
		return ErrNoOutputChannels
	}
	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}
			target := outs[i%len(outs)]
			select {
			case target <- data:
				i++
			case <-ctx.Done():
				return nil
			}
		}
	}
}
