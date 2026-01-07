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
		chCopy := ch
		go func() {
			defer wg.Done()
			for val := range chCopy {
				if strings.Contains(val, "no multiplexer") {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- val:
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
