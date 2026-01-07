package handlers

import (
	"context"
	"strings"
	"sync"
)

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var wg sync.WaitGroup

	for _, input := range inputs {
		inputCh := input

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, ok := <-inputCh:
					if !ok {
						return
					}

					if strings.Contains(value, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- value:
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
