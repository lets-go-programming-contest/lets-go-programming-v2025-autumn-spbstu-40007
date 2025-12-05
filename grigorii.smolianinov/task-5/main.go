package main

import (
	"context"
	"log"
	"time"

	"github.com/Smolyaninoff/GoLang/pkg/conveyer"
	"github.com/Smolyaninoff/GoLang/pkg/handlers"
)

func main() {
	c := conveyer.New(5)

	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input_A", "output_A")
	c.RegisterSeparator(handlers.SeparatorFunc, "output_A", []string{"output_B1", "output_B2"})
	c.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"output_B1", "output_B2"}, "final_output")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		log.Println("Конвейер запущен...")
		if err := c.Run(ctx); err != nil {
			log.Printf("Конвейер завершился с ошибкой: %v", err)
		} else {
			log.Println("Конвейер завершился успешно по истечении времени.")
		}
	}()

	log.Println("Начинаем отправку данных в 'input_A'...")

	dataToSend := []string{
		"Hello",
		"World",
		"Test no decorator",
		"Data 4",
		"Data 5",
	}

	for i, data := range dataToSend {
		if i == 2 {
			log.Printf("Отправка данных '%s', что вызовет ошибку...", data)
		}

		if err := c.Send("input_A", data); err != nil {
			log.Printf("Ошибка при отправке '%s': %v (Возможно, канал уже закрыт)", data, err)
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	<-ctx.Done()
	log.Println("Горутина main ждет завершения конвейера...")

	log.Println("Попытка получения оставшихся данных из 'final_output'...")
	val, err := c.Recv("final_output")
	if err != nil {
		log.Printf("Recv ошибка: %v", err)
	} else {
		log.Printf("Recv результат: %s", val)
	}

	log.Println("Демонстрация завершена.")
}
