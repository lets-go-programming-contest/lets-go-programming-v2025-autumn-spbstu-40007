package handlers

import "context"

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	idx := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case s, ok := <-input:
			if !ok {
				for _, ch := range outputs {
					close(ch)
				}
				return nil
			}

			out := outputs[idx]
			idx = (idx + 1) % len(outputs)

			select {
			case out <- s:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
