package main

import (
	"flag"
	"log"

	"github.com/UwUshkin/task-3/internal/config"
	"github.com/UwUshkin/task-3/internal/processor"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to the YAML configuration file")

	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Fatal error loading config file '%s': %v", configPath, err)
	}


	if err := processor.ProcessAndSave(cfg.InputFile, cfg.OutputFile); err != nil {
		log.Fatalf("Fatal error during data processing: %v", err)
	}

}
