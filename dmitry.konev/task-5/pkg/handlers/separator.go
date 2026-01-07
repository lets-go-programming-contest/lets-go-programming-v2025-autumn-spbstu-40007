package handlers

import "context"

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	i := 0
	count := len(outputs)
	for val := range input {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case outputs[i%count] <- val:
			i++
		}
	}
	return nil
}
