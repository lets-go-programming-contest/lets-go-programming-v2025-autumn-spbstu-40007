package handlers

import (
	"context"
	"fmt"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := counter % len(outputs)
			select {
			case <-ctx.Done():
				return nil
			case outputs[idx] <- data:
				counter++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	// Создаем каналы для отслеживания активных входов
	activeInputs := inputs

	for len(activeInputs) > 0 {
		for i := 0; i < len(activeInputs); i++ {
			select {
			case <-ctx.Done():
				return nil
			case data, ok := <-activeInputs[i]:
				if !ok {
					// Удаляем закрытый канал
					activeInputs = append(activeInputs[:i], activeInputs[i+1:]...)
					i--
					continue
				}

				if strings.Contains(data, "no multiplexer") {
					continue
				}

				select {
				case <-ctx.Done():
					return nil
				case output <- data:
				}
			default:
			}
		}

		// Если все каналы закрыты
		if len(activeInputs) == 0 {
			return nil
		}
	}

	return nil
}
