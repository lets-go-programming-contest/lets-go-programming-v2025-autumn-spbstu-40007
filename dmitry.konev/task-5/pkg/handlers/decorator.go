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
			return ctx.Err()

		case s, ok := <-input:
			if !ok {
				close(output)
				return nil
			}

			if strings.Contains(s, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(s, "decorated: ") {
				s = "decorated: " + s
			}

			select {
			case output <- s:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
