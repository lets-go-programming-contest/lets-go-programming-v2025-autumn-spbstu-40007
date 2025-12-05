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
				// Если входной канал закрыт, завершаем работу
				return nil
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

	var counter int64 = 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				// Если входной канал закрыт
				return nil
			}

			// Выбираем выходной канал по порядку
			idx := atomic.AddInt64(&counter, 1) - 1
			outputIdx := int(idx) % len(outputs)

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

	// Для каждого входного канала создаем горутину
	type result struct {
		data string
		ok   bool
	}

	merged := make(chan result, len(inputs)*10)

	// Запускаем горутины для каждого входного канала
	for _, input := range inputs {
		go func(in chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}

					select {
					case <-ctx.Done():
						return
					case merged <- result{data: data, ok: ok}:
					}
				}
			}
		}(input)
	}

	// Обрабатываем объединенные данные
	for {
		select {
		case <-ctx.Done():
			return nil
		case res, ok := <-merged:
			if !ok {
				return nil
			}

			// Фильтрация данных с подстрокой "no multiplexer"
			if strings.Contains(res.data, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- res.data:
			}
		}
	}
}
