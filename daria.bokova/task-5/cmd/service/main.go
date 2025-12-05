package main

import (
	"context"
	"fmt"
	"time"

	"conveyer/pkg/conveyer"
	"conveyer/pkg/handlers"
)

func main() {
	// Простой пример использования конвейера
	c := conveyer.New(100)

	// Регистрируем обработчики
	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input1", "output1")
	c.RegisterSeparator(handlers.SeparatorFunc, "output1", []string{"out1", "out2"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Запускаем конвейер
	go func() {
		if err := c.Run(ctx); err != nil {
			fmt.Printf("Conveyer stopped: %v\n", err)
		}
	}()

	// Даем время на запуск
	time.Sleep(100 * time.Millisecond)

	// Отправляем данные
	c.Send("input1", "Test message 1")
	c.Send("input1", "Test message 2")

	// Получаем результаты
	for i := 0; i < 2; i++ {
		if data, err := c.Recv("out1"); err == nil {
			fmt.Printf("out1: %s\n", data)
		}
		if data, err := c.Recv("out2"); err == nil {
			fmt.Printf("out2: %s\n", data)
		}
	}
}
