package main

import (
	"fmt"
	"log"

	"github.com/Mart22052006/task-8/pkg/config"
)

func main() {
	configuration, loadErr := config.Load()
	if loadErr != nil {
		log.Fatalf("failed to load configuration: %v", loadErr)
	}

	fmt.Printf("%s %s", configuration.Environment, configuration.LogLevel)
}
