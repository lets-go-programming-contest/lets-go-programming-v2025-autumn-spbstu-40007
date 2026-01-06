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

	for _, ch := range inputs {
		ch := ch
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(v, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- v:
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}