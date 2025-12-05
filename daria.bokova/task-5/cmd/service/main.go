package main

import (
	"context"
	"fmt"
	"time"

	"conveyer/pkg/conveyer"
	"conveyer/pkg/handlers"
)

func main() {
	// Тест 1: Базовый тест
	fmt.Println("=== Test 1: Basic Test ===")
	testBasic()

	// Тест 2: Тест с ошибкой декоратора
	fmt.Println("\n=== Test 2: Decorator Error Test ===")
	testDecoratorError()

	// Тест 3: Тест сепаратора
	fmt.Println("\n=== Test 3: Separator Test ===")
	testSeparator()
}

func testBasic() {
	c := conveyer.New(10)

	// Регистрируем обработчики
	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "decorated")
	c.RegisterSeparator(handlers.SeparatorFunc, "decorated", []string{"out1", "out2"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Запускаем конвейер
	go func() {
		if err := c.Run(ctx); err != nil {
			fmt.Printf("Conveyer stopped: %v\n", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	// Отправляем данные
	c.Send("input", "Hello")
	c.Send("input", "World")

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

func testDecoratorError() {
	c := conveyer.New(10)

	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "output")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Запускаем конвейер
	go func() {
		if err := c.Run(ctx); err != nil {
			fmt.Printf("Conveyer stopped with error: %v\n", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	// Отправляем данные с ошибкой
	err := c.Send("input", "This has no decorator in it")
	if err != nil {
		fmt.Printf("Send error: %v\n", err)
	}

	time.Sleep(100 * time.Millisecond)
}

func testSeparator() {
	c := conveyer.New(10)

	c.RegisterSeparator(handlers.SeparatorFunc, "input", []string{"out1", "out2", "out3"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Запускаем конвейер
	go func() {
		if err := c.Run(ctx); err != nil {
			fmt.Printf("Conveyer stopped: %v\n", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	// Отправляем несколько сообщений
	for i := 1; i <= 5; i++ {
		c.Send("input", fmt.Sprintf("Message %d", i))
	}

	// Получаем результаты
	for i := 1; i <= 5; i++ {
		if data, err := c.Recv("out1"); err == nil && data != "" {
			fmt.Printf("out1: %s\n", data)
		}
		if data, err := c.Recv("out2"); err == nil && data != "" {
			fmt.Printf("out2: %s\n", data)
		}
		if data, err := c.Recv("out3"); err == nil && data != "" {
			fmt.Printf("out3: %s\n", data)
		}
	}
}
