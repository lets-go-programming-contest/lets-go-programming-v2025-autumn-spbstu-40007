package main

import (
	"fmt"
	"log"

	"github.com/Mart22052006/task-8/pkg/config"
)

func main() {
	loadedConfig, parseErr := config.ParseConfiguration(config.EmbeddedConfig)
	if parseErr != nil {
		log.Fatalf("Failed to load config: %v", parseErr)
	}

	fmt.Printf("%s %s\n", loadedConfig.Environment, loadedConfig.LogLevel)
}
