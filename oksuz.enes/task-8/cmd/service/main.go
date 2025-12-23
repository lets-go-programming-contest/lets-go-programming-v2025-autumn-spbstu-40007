package main

import (
	"fmt"

	"oksuz.enes/task-8/config"
)

func main() {
	fmt.Print(config.Environment, " ", config.LogLevel)
}
