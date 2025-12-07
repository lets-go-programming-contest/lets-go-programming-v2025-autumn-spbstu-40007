package main

import (
	"fmt"

	"github.com/tntkatz/task-8/pkg/config"
)

func main() {
	cfg := config.Load()

	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
