package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ksuah/task-5/pkg/conveyer"
	"github.com/ksuah/task-5/pkg/handlers"
)

func main() {
	rootCtx, rootCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer rootCancel()

	c := conveyer.New(5)

	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "in_1", "mid_1")
	c.RegisterSeparator(handlers.SeparatorFunc, "mid_1", []string{"out_A", "out_B"})
	c.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"out_A", "out_B"}, "final_out")

	doneChan := make(chan error)

	go func() {
		fmt.Println("Конвейер запущен")
		err := c.Run(rootCtx)
		doneChan <- err
	}()

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\nОтправка данных в in_1")

	messages := []string{"Hello", "World", "GoLang", "no multiplexer data", "test"}
	for i, msg := range messages {
		fmt.Printf("SEND: %s\n", msg)
		c.Send("in_1", fmt.Sprintf("%d:%s", i+1, msg))
	}

	fmt.Println("\nПолучение данных из final_out")

	for i := 0; i < len(messages); i++ {
		data, err := c.Recv("final_out")
		if err != nil {
			fmt.Printf("RECV ERROR: %v\n", err)
			break
		}
		fmt.Printf("RECV: %s\n", data)
	}

	fmt.Println("\nТест ошибки (Decorator)")

	err := c.Send("in_1", "Fatal error: no decorator detected")
	if err == nil {
		fmt.Println("SEND: Fatal error: no decorator detected (Ожидается ошибка в обработчике)")
	}

	finalErr := <-doneChan

	if finalErr != nil && finalErr != context.Canceled && finalErr != context.DeadlineExceeded {
		fmt.Printf("Конвейер завершился корректно по ошибке: %v\n", finalErr)
	} else {
		fmt.Printf("Конвейер завершился НЕ из-за ошибки обработчика: %v\n", finalErr)
	}

	_, err = c.Recv("final_out")
	if err != nil {
		fmt.Printf("RECV ERROR (ОЖИДАЕМАЯ ОСТАНОВКА): %v\n", err)
	} else {
		fmt.Println("ОШИБКА: Recv получил данные, хотя конвейер должен быть остановлен.")
	}

	fmt.Println("\nКонвейер завершил работу")
}
