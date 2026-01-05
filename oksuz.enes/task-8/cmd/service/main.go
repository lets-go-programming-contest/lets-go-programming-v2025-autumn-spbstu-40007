package main

import (
	"fmt"

	"oksuz.enes/task-8/config"
)

func main() {
	cfg := config.Load()
	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
