package main

import (
	"fmt"

	"github.com/falsefeelings/task-8/pkg/config"
)

func main() {
	config := config.Load()

	fmt.Printf("%s %s", config.Environment, config.LogLevel)

}
