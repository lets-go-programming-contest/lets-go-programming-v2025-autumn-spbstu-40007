package main

import (
	"fmt"

	"oksuz.enes/task-8/config"
)

func main() {
	cfg := config.NewConfig()
	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
