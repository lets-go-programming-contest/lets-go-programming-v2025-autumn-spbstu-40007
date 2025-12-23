package main

import (
	"fmt"

	"oksuz.enes/task-8/config"
)

func main() {
	fmt.Println(config.Cfg.Environment, config.Cfg.LogLevel)
}
