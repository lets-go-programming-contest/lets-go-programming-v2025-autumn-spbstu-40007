package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

func safelyCloseChannel(chnl chan string) {
	defer func() {
		if r := recover(); r != nil {
			
		}
	}()
	if chnl != nil {
		close(chnl)
	}
}

func readInputToTransfer(ctx context.Context, input chan string, transfer chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-input:
			if !ok {
				return
			}

			select {
			case transfer <- data:
			case <-ctx.Done():
				return
			}
		}
	}
}

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	defer safelyCloseChannel(output)

	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return errors.New("data error: can't be decorated because it contains 'no decorator'")
			}

			processedData := data
			if !strings.HasPrefix(data, prefix) {
				processedData = prefix + data
			}

			select {
			case output <- processedData:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer func() {
		closedChans := make(map[chan string]struct{})
		for _, output := range outputs {
			if _, ok := closedChans[output]; !ok {
				safelyCloseChannel(output)
				closedChans[output] = struct{}{}
			}
		}
	}()

	if len(outputs) == 0 {
		for range input {
			select {
			case <-ctx.Done():
				return nil
			default:
			}
		}
		return nil
	}

	var counter uint64
	var mu sync.Mutex
	numOutputs := uint64(len(outputs))

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			mu.Lock()
			index := counter % numOutputs
			counter++
			mu.Unlock()

			output := outputs[index]

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	defer safelyCloseChannel(output)

	var wg sync.WaitGroup
	transfer := make(chan string)

	for _, input := range inputs {
		wg.Add(1)
		go readInputToTransfer(ctx, input, transfer, &wg)
	}

	go func() {
		wg.Wait()
		safelyCloseChannel(transfer)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-transfer:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no multiplexer") {
				continue
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}