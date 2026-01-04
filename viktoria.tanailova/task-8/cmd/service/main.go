package main

import (
	"fmt"

	"task-8/config"
)

func main() {
	cfg := config.Load()
	fmt.Println(cfg.Environment, cfg.LogLevel)
}
