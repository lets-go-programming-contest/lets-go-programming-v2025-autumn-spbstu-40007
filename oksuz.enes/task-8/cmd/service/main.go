package main

import (
	"fmt"

	"oksuz.enes/task-8/config"
)

func main() {
	fmt.Printf("%s %s", config.Environment, config.LogLevel)
}
