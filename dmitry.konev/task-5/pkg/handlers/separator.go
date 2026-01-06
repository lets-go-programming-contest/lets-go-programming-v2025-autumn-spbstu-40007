package handlers

import "context"

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	index := 0
	count := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case value, ok := <-input:
			if !ok {
				return nil
			}

			outputs[index%count] <- value
			index++
		}
	}
}
