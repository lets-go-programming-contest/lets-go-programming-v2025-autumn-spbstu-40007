package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
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
			select {
			case <-ctx.Done():
				return nil
			case output <- val:
			}
		}
	}
}
