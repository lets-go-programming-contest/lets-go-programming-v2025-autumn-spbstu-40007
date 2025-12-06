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

func PrefixDecoratorFunc(ctx context.Context, in, out chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}
			if strings.Contains(data, "no decorator") {
				return ErrCannotDecorate
			}
			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
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
	var group sync.WaitGroup
	group.Add(len(ins))

	for _, input := range ins {
		go func(ch <-chan string) {
			defer group.Done()
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
					case out <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(input)
	}

	group.Wait()
	return nil
}

func SeparatorFunc(ctx context.Context, in chan string, outs []chan string) error {
	if len(outs) == 0 {
		return ErrNoOutputChannels
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}
			select {
			case outs[index%len(outs)] <- data:
				index++
			case <-ctx.Done():
				return nil
			}
		}
	}
}
