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
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(inputs))

	for _, inputCh := range inputs {
		channel := inputCh

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(value, "no multiplexer") {
						continue
					}

					output <- value
				}
			}
		}()
	}

	waitGroup.Wait()
	return nil
}