package main

import (
	"flag"
	"log"

	"github.com/Mart22052006/task-3/internal/configuration"
	"github.com/Mart22052006/task-3/internal/currencyprocessor"
)

func main() {
	configurationPath := flag.String(
		"config",
		"./configs/config.yaml",
		"Path to YAML configuration file containing input and output file paths",
	)

	flag.Parse()

	applicationConfig, configError := configuration.LoadApplicationConfig(*configurationPath)
	if configError != nil {
		log.Fatalf("Failed to load configuration: %v", configError)
	}

	processingError := currencyprocessor.ExecuteProcessing(applicationConfig)
	if processingError != nil {
		log.Fatalf("Failed to process currencies: %v", processingError)
	}
}
