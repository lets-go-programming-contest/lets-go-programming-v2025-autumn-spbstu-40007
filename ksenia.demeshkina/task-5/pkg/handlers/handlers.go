package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, inputChannel chan string, outputChannel chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}

			newData := data

			if !strings.HasPrefix(data, "decorated: ") {
				newData = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return nil
			case outputChannel <- newData:
			}
		}
	}
}

func handleInputChannel(
	ctx context.Context,
	inputChannel chan string,
	transferChannel chan string,
	waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case val, ok := <-inputChannel:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return
			case transferChannel <- val:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	var waitGroup sync.WaitGroup

	transferChannel := make(chan string)

	for _, inputChan := range inputChannels {
		waitGroup.Add(1)

		go handleInputChannel(ctx, inputChan, transferChannel, &waitGroup)
	}

	go func() {
		waitGroup.Wait()
		close(transferChannel)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-transferChannel:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case outputChannel <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, inputChannel chan string, outputChannels []chan string) error {
	var index int

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			targetIdx := index % len(outputChannels)

			index++

			select {
			case <-ctx.Done():
				return nil
			case outputChannels[targetIdx] <- data:
			}
		}
	}
}
