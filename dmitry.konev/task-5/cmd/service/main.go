package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DichSwitch/task-5/pkg/conveyer"
	"github.com/DichSwitch/task-5/pkg/handlers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := conveyer.New(16)

	p.RegisterDecorator(handlers.PrefixDecoratorFunc, "in", "after-decorator")
	p.RegisterSeparator(handlers.SeparatorFunc, "after-decorator", []string{"s1", "s2"})
	p.RegisterMultiplexer(handlers.MultiplexerFunc, []string{"s1", "s2"}, "out")

	runErr := make(chan error, 1)
	go func() {
		runErr <- p.Run(ctx)
	}()

	inputs := []string{
		"hello",
		"world",
		"no multiplexer example",
		"foo",
		"no decorator should fail",
	}

	for _, s := range inputs {
		err := p.Send("in", s)
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}

	time.Sleep(300 * time.Millisecond)

	for {
		v, err := p.Recv("out")
		if err != nil {
			fmt.Println("Recv error:", err)
			break
		}
		if v == "undefined" {
			fmt.Println("Output channel closed and drained.")
			break
		}
		fmt.Println("OUT:", v)
	}

	cancel()

	select {
	case err := <-runErr:
		if err != nil {
			fmt.Println("Pipeline finished with error:", err)
		} else {
			fmt.Println("Pipeline finished successfully")
		}
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout waiting for pipeline to finish")
	}
}
