package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrDecorationSkipped = errors.New("can't be decorated")

const (
	skipDecorationKeyword = "no decorator"
	decorationPrefix      = "decorated: "
	skipMuxKeyword        = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, inStream chan string, outStream chan string) error {
	defer close(outStream)

	for {
		select {
		case <-ctx.Done():
			return nil

		case dataChunk, ok := <-inStream:
			if !ok {
				return nil
			}

			if strings.Contains(dataChunk, skipDecorationKeyword) {
				return ErrDecorationSkipped
			}

			if !strings.HasPrefix(dataChunk, decorationPrefix) {
				dataChunk = decorationPrefix + dataChunk
			}

			select {
			case <-ctx.Done():
				return nil
			case outStream <- dataChunk:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inStreams []chan string, outStream chan string) error {
	defer close(outStream)

	if len(inStreams) == 0 {
		return nil
	}

	for _, sourceStream := range inStreams {
		localSourceStream := sourceStream

		go func(chLocal chan string) {
			for {
				select {
				case <-ctx.Done():
					return

				case dataChunk, ok := <-chLocal:
					if !ok {
						return
					}

					if strings.Contains(dataChunk, skipMuxKeyword) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case outStream <- dataChunk:
					}
				}
			}
		}(localSourceStream)
	}

	<-ctx.Done()

	return nil
}

func SeparatorFunc(ctx context.Context, inStream chan string, outStreams []chan string) error {
	defer func() {
		for _, outCh := range outStreams {
			close(outCh)
		}
	}()

	sinkCount := len(outStreams)
	dispatchIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case dataChunk, ok := <-inStream:
			if !ok {
				return nil
			}

			if sinkCount == 0 {
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case outStreams[dispatchIndex%sinkCount] <- dataChunk:
				dispatchIndex++
			}
		}
	}
}
