package handlers

import (
	"context"
	"strings"
	"sync"
)

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

	wg.Wait()
	return nil
}