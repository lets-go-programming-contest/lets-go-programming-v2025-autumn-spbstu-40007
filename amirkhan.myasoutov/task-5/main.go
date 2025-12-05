package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ami0-0/task-5/pkg/conveyer"
	"github.com/ami0-0/task-5/pkg/handlers"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	fmt.Println("Инициализация конвейера")
	c := conveyer.New(5)

	c.RegisterSeparator(handlers.SeparatorFunc, "in", []string{"sep1", "sep2"})
	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "sep1", "dec1")
	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "sep2", "dec2")
	c.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"dec1", "dec2"}, "out")

	go func() {
		err := c.Run(ctx)
		if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
			fmt.Printf("\nКонвейер завершен с ошибкой: %v\n", err)
		} else {
			fmt.Println("\nКонвейер завершил работу.")
		}
	}()
	fmt.Println("Конвейер запущен.")

	fmt.Println("\n--- Отправка данных ---")
	c.Send("in", "Hello A")
	c.Send("in", "Hello B")
	c.Send("in", "no decorator example")
	c.Send("in", "Test C")

	time.Sleep(300 * time.Millisecond)

	fmt.Println("\n--- Получение данных ---")
	for i := 0; i < 5; i++ {
		data, err := c.Recv("out")
		if err != nil {
			fmt.Printf("Recv ошибка: %v\n", err)
			break
		}
		if data == "undefined" {
			fmt.Println("Recv: Канал закрыт (undefined).")
			break
		}
		fmt.Printf("Recv [%d]: %s\n", i+1, data)
	}

	time.Sleep(1 * time.Second)
}
