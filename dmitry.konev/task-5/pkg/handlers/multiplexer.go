package handlers

import (
	"context"
	"strings"
	"sync"
)

func CombineChannels(ctx context.Context, inputs []chan string, out chan string) error {
	var wg sync.WaitGroup
	for _, ch := range inputs {
		c := ch
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-c:
					if !ok {
						return
					}
					if strings.Contains(val, "no multiplexer") {
						continue
					}
					select {
					case <-ctx.Done():
						return
					case out <- val:
					}
				}
			}
		}()
	}
	wg.Wait()
	return nil
}
