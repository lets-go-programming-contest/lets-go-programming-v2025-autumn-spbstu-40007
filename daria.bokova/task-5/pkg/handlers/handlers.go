package handlers

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
)

// PrefixDecoratorFunc - модификатор данных
func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				// Если входной канал закрыт, проверяем контекст
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			}

			// Проверяем наличие подстроки "no decorator"
			if strings.Contains(data, "no decorator") {
				return errors.New("can't be decorated")
			}

			// Добавляем префикс, если его еще нет
			var result string
			if strings.HasPrefix(data, prefix) {
				result = data
			} else {
				result = prefix + data
			}

			// Отправляем результат
			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- result:
			}
		}
	}
}

// SeparatorFunc - сепаратор по порядковому номеру
func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var counter int64 = 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				// Если входной канал закрыт
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			}

			// Выбираем выходной канал по порядку
			idx := atomic.AddInt64(&counter, 1) - 1
			outputIdx := int(idx) % len(outputs)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[outputIdx] <- data:
			}
		}
	}
}

// MultiplexerFunc - мультиплексор с фильтрацией
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	// Создаем канал для объединения входов
	merged := make(chan string, 100)

	// Запускаем горутины для каждого входного канала
	done := make(chan struct{})
	defer close(done)

	for _, input := range inputs {
		go func(in chan string) {
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}

					select {
					case <-done:
						return
					case <-ctx.Done():
						return
					case merged <- data:
					}
				}
			}
		}(input)
	}

	// Обрабатываем объединенные данные
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-merged:
			if !ok {
				return nil
			}

			// Фильтрация данных с подстрокой "no multiplexer"
			if strings.Contains(data, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- data:
			}
		}
	}
}
