package handlers

import (
	"context"
	"strings"
)

func MultiplexerFunc(
	ctx context.Context,
	inputChans []chan string,
	outputChan chan string,
) error {
	openChans := len(inputChans)
	doneCh := make(chan struct{}, len(inputChans))

	for _, inputChan := range inputChans {
		go func(in chan string) {
			for {
				select {
				case <-ctx.Done():
					return

				case val, ok := <-in:
					if !ok {
						doneCh <- struct{}{}
						return
					}

					if strings.Contains(val, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case outputChan <- val:
					}
				}
			}
		}(inputChan)
	}

	for openChans > 0 {
		select {
		case <-ctx.Done():
			return nil
		case <-doneCh:
			openChans--
		}
	}

	return nil
}