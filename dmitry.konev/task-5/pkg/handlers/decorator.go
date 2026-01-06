package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(value, "decorated: ") {
				value = "decorated: " + value
			}

			output <- value
		}
	}
}