package main

import (
	"fmt"

	"github.com/tntkatz/task-8/pkg/config"
)

func main() {
	cfg := config.Load()

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
