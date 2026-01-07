package handlers

import (
	"context"
	"strings"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	for _, in := range inputs {
		ch := in
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-ch:
					if !ok {
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
		}()
	}

	<-ctx.Done()
	return ctx.Err()
}
