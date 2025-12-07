package handlers

import (
	"context"
)

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {

	open := len(inputs)
	closed := make([]bool, len(inputs))

	for open > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		for i, ch := range inputs {
			if closed[i] {
				continue
			}

			select {
			case s, ok := <-ch:
				if !ok {
					closed[i] = true
					open--
					continue
				}

				select {
				case output <- s:
				case <-ctx.Done():
					return ctx.Err()
				}

			default:
			}
		}
	}

	close(output)
	return nil
}
