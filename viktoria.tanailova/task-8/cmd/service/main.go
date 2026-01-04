package main

import (
	"fmt"
	"task-8/config"
)

func main() {
	cfg := config.Load()
	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
