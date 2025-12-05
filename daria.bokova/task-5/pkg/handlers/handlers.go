package handlers

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return errors.New("can't be decorated")
			}

			var result string
			if strings.HasPrefix(data, prefix) {
				result = data
			} else {
				result = prefix + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- result:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var idx int64 = 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			// Используем атомарный счетчик
			outputIdx := int(atomic.AddInt64(&idx, 1)-1) % len(outputs)

			select {
			case <-ctx.Done():
				return nil
			case outputs[outputIdx] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	// Простой мультиплексор - читаем из всех входов по очереди
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Проверяем все входы
			for _, input := range inputs {
				select {
				case <-ctx.Done():
					return nil
				case data, ok := <-input:
					if !ok {
						continue
					}

					// Фильтрация
					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return nil
					case output <- data:
					}
				default:
					// Нет данных в этом канале
				}
			}
		}
	}
}
