package handlers

import "context"

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	index := 0
	count := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-input:
			if !ok {
				return nil
			}

			ch := outputs[index%count]
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ch <- v:
			}
			index++
		}
	}
}