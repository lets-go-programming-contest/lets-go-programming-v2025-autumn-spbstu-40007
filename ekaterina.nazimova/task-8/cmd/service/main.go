package main

import (
	"fmt"
	"github.com/UwUshkin/task-8/pkg/config"
)

func main() {
	cfg := config.GetConfig()
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
