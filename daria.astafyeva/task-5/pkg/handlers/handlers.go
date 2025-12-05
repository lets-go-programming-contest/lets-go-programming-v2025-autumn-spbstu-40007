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
	skipMux   = "no multiplexer"
	skipDec   = "no decorator"
	decPrefix = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, inCh chan string, outCh chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, open := <-inCh:
			if !open {
				return nil
			}

			if strings.Contains(data, skipDec) {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(data, decPrefix) {
				data = decPrefix + data
			}

			select {
			case outCh <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inChs []chan string, outCh chan string) error {
	var wg sync.WaitGroup

	for _, in := range inChs {
		wg.Add(1)
		go func(ch <-chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, open := <-ch:
					if !open {
						return
					}

					if strings.Contains(data, skipMux) {
						continue
					}

					select {
					case outCh <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(in)
	}

	wg.Wait()
	return nil
}

func SeparatorFunc(ctx context.Context, inCh chan string, outChs []chan string) error {
	if len(outChs) == 0 {
		return ErrNoOutputChannels
	}

	var pos int
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, open := <-inCh:
			if !open {
				return nil
			}

			target := outChs[pos%len(outChs)]
			select {
			case target <- data:
				pos++
			case <-ctx.Done():
				return nil
			}
		}
	}
}
