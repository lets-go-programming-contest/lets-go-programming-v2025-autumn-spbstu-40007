package handlers

import (
	"context"
	"reflect"
	"strings"
)

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {

	cases := make([]reflect.SelectCase, len(inputs))
	open := len(inputs)

	for i, ch := range inputs {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	for open > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		idx, val, ok := reflect.Select(cases)
		if !ok {
			cases[idx].Chan = reflect.ValueOf(nil)
			open--
			continue
		}

		s := val.String()

		if strings.Contains(s, "no multiplexer") {
			continue
		}

		select {
		case output <- s:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	close(output)
	return nil
}
