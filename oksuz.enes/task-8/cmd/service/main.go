package main

import (
	"fmt"
	"task-8/config"
)

func main() {
	fmt.Printf("%s %s\n", config.Environment, config.LogLevel)
}
