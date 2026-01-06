package handlers

import "context"

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return nil
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case val, ok := <-input:
			if !ok {
				return nil
			}

			out := outputs[index%len(outputs)]
			index++

			select {
			case <-ctx.Done():
				return nil
			case out <- val:
			}
		}
	}
}

