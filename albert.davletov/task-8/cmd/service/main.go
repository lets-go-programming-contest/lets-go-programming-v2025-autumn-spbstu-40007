package main

import (
	"fmt"

	"task-8/config"
)

func main() {
	c := config.GetConfig()

	fmt.Print(c.Environment, " ", c.LogLevel)
}
