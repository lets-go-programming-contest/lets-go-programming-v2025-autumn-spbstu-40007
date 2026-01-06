package handlers

import "context"

func SeparatorFunc(
	ctx context.Context,
	inputChan chan string,
	outputChans []chan string,
) error {
	if len(outputChans) == 0 {
		return nil
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case val, ok := <-inputChan:
			if !ok {
				return nil
			}

			targetChan := outputChans[index%len(outputChans)]
			index++

			select {
			case <-ctx.Done():
				return nil
			case targetChan <- val:
			}
		}
	}
}