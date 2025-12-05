package main

import (
	"context"
	"fmt"
	"time"

	"conveyer/pkg/conveyer"
	"conveyer/pkg/handlers"
)

func main() {
	// Создаем конвейер
	c := conveyer.New(100)

	// Регистрируем обработчики
	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input1", "decorated1")
	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input2", "decorated2")

	c.RegisterSeparator(handlers.SeparatorFunc, "decorated1", []string{"sep_out1", "sep_out2"})

	c.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"sep_out1", "sep_out2"}, "final_output")

	// Запускаем конвейер в отдельной горутине
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := c.Run(ctx); err != nil {
			fmt.Printf("Conveyer stopped: %v\n", err)
		}
	}()

	// Даем время на запуск
	time.Sleep(100 * time.Millisecond)

	// Отправляем данные
	c.Send("input1", "Hello World")
	c.Send("input2", "Test message")
	c.Send("input1", "Another message")
	c.Send("input2", "message with no multiplexer")

	// Получаем результаты
	for i := 0; i < 3; i++ {
		if data, err := c.Recv("final_output"); err == nil {
			fmt.Printf("Received: %s\n", data)
		}
	}

	// Тест с ошибкой
	err := c.Send("input1", "This has no decorator in it")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Даем время на обработку
	time.Sleep(500 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)

	// Пробуем получить данные из закрытого канала
	if data, err := c.Recv("final_output"); err == nil {
		fmt.Printf("From closed channel: %s\n", data)
	}
}
