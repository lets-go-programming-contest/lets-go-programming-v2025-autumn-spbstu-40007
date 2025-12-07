package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCannotDecorate = errors.New("can't be decorated")

func safeClose(ch chan string) {
	defer func() { _ = recover() }()
	close(ch)
}

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	defer safeClose(output)

	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():

			return nil
		case v, ok := <-input:
			if !ok {

				return nil
			}

			if strings.Contains(v, "no decorator") {

				return ErrCannotDecorate
			}

			if !strings.HasPrefix(v, prefix) {
				v = prefix + v
			}

			select {
			case output <- v:
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
	defer safeClose(output)

	for {
		allClosed := true

		for _, in := range inputs {
			select {
			case <-ctx.Done():

				return nil
			default:
			}

			select {
			case v, ok := <-in:
				if !ok {
					continue
				}
				allClosed = false

				if strings.Contains(v, "no multiplexer") {
					continue
				}

				select {
				case output <- v:
				case <-ctx.Done():

					return nil
				}
			default:
			}
		}

		if allClosed {

			return nil
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {

		return nil
	}

	defer func() {
		seen := make(map[chan string]struct{})
		for _, ch := range outputs {
			if _, ok := seen[ch]; ok {
				continue
			}
			seen[ch] = struct{}{}
			safeClose(ch)
		}
	}()

	idx := 0

	for {
		select {
		case <-ctx.Done():

			return nil
		case v, ok := <-input:
			if !ok {

				return nil
			}

			out := outputs[idx%len(outputs)]
			idx++

			select {
			case out <- v:
			case <-ctx.Done():

				return nil
			}
		}
	}
}
