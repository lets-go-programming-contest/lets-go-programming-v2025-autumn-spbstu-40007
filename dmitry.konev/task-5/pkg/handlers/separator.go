package handlers

import "context"

func SplitChannel(ctx context.Context, in chan string, outs []chan string) error {
	if len(outs) == 0 {
		return nil
	}
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-in:
			if !ok {
				return nil
			}
			select {
			case <-ctx.Done():
				return nil
			case outs[idx%len(outs)] <- val:
				idx++
			}
		}
	}
}
