package handlers

import (
	"context"
	"errors"
	"strings"
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}

			output <- val
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	idx := 0
	count := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return nil

		case val, ok := <-input:
			if !ok {
				return nil
			}

			outputs[idx%count] <- val
			idx++
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	open := len(inputs)

	for open > 0 {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		for i := 0; i < len(inputs); i++ {
			select {
			case val, ok := <-inputs[i]:
				if !ok {
					inputs[i] = nil
					open--
					continue
				}

				if strings.Contains(val, "no multiplexer") {
					continue
				}

				output <- val
			default:
			}
		}
	}
	return nil
}