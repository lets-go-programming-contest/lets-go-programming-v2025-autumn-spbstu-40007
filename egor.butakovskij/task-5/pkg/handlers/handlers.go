package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var errPrefixDecoratorFuncCantBeDecorated = errors.New("can't be decorated")

func safelyCloseChannel(chnl chan string) {
	defer func() {
		if r := recover(); r != nil {
			_ = r
		}
	}()
	close(chnl)
}

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	defer safelyCloseChannel(output)

	prefix := "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return errPrefixDecoratorFuncCantBeDecorated
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

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var wgrp sync.WaitGroup

	transfer := make(chan string)

	defer safelyCloseChannel(output)

	for _, input := range inputs {
		wgrp.Add(1)
		go readInputToTransfer(ctx, input, transfer, &wgrp) //nolint:wsl
	}

	go func() {
		wgrp.Wait()
		close(transfer)
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

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer (func() {
		closedChans := make(map[chan string]struct{})
		for _, output := range outputs {
			if _, ok := closedChans[output]; ok {
				continue
			}
			closedChans[output] = struct{}{}
			safelyCloseChannel(output)
		}
	})()

	var counter int

	if len(outputs) == 0 {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			targetIndex := counter % len(outputs)
			targetChan := outputs[targetIndex]

			select {
			case targetChan <- data:
				counter++
			case <-ctx.Done():
				return nil
			}
		}
	}
}
