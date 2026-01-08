package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantDecorate = errors.New("can't be decorated")

func AddPrefix(ctx context.Context, in chan string, out chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-in:
			if !ok {
				return nil
			}
			if strings.Contains(val, "no decorator") {
				return ErrCantDecorate
			}
			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}
			select {
			case <-ctx.Done():
				return nil
			case out <- val:
			}
		}
	}
}
