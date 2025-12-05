package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bdshka/go-conveyer/pkg/conveyer"
	"github.com/bdshka/go-conveyer/pkg/handlers"
)

func main() {
	c := conveyer.New(10)

	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input1", "decorated")
	c.RegisterSeparator(handlers.SeparatorFunc, "decorated", []string{"out1", "out2"})
	c.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"out1", "out2"}, "final")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go func() {
		if err := c.Run(ctx); err != nil {
			fmt.Printf("Conveyer error: %v\n", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	fmt.Println("Sending data...")
	c.Send("input1", "hello")
	c.Send("input1", "world")
	c.Send("input1", "test without decorator")

	fmt.Println("\nReceiving data...")
	for i := 0; i < 3; i++ {
		val, err := c.Recv("final")
		if err != nil {
			fmt.Printf("Error receiving: %v\n", err)
		} else {
			fmt.Printf("Received: %s\n", val)
		}
	}

	fmt.Println("\nSending data that should cause error...")
	c.Send("input1", "no decorator test")

	time.Sleep(1 * time.Second)
	fmt.Println("\nDone!")
}
