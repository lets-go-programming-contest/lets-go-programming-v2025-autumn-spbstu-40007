package handlers

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync/atomic"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				// Входной канал закрыт
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
				return ctx.Err()
			case output <- result:
				// Успешно отправлено
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				// Входной канал закрыт
				return nil
			}

			// Выбираем выходной канал по порядку
			idx := atomic.AddUint64(&counter, 1) - 1
			outputIdx := int(idx) % len(outputs)

			// Отправляем в выбранный канал
			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[outputIdx] <- data:
				// Успешно отправлено
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	// Используем select для чтения из всех каналов
	cases := make([]reflect.SelectCase, len(inputs)+1)
	for i, ch := range inputs {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}
	// Добавляем case для контекста
	cases[len(inputs)] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	}

	for {
		chosen, value, ok := reflect.Select(cases)

		if chosen == len(inputs) {
			// Контекст завершен
			return ctx.Err()
		}

		if !ok {
			// Канал закрыт
			// Заменяем его на nil, чтобы больше не выбирать
			cases[chosen].Chan = reflect.ValueOf(nil)
			continue
		}

		data := value.String()

		// Фильтрация
		if strings.Contains(data, "no multiplexer") {
			continue
		}

		// Отправляем в выходной канал
		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- data:
			// Успешно
		}
	}
}
