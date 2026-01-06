package handlers

import (
	"context"
	"strings"
)

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	open := len(inputs)
	done := make(chan struct{}, len(inputs))

	for _, ch := range inputs {
		go func(ch chan string) {
			for {
				select {
				case <-ctx.Done():
					return

				case val, ok := <-ch:
					if !ok {
						done <- struct{}{}
						return
					}

					if strings.Contains(val, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- val:
					}
				}
			}
		}(ch)
	}

	for open > 0 {
		select {
		case <-ctx.Done():
			return nil
		case <-done:
			open--
		}
	}

	return nil
}
