package main

import (
	"fmt"
	"task-8/config" // Импорт через имя корневого модуля
)

func main() {
	cfg := config.Load()
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
