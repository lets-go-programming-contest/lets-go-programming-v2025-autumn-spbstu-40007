package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	defer close(output)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(val, "no decorator") {
				return ErrCantBeDecorated
			}
			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}
			output <- val
		}
	}
}
