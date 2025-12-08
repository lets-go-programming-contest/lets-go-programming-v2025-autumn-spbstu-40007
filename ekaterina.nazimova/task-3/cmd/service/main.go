package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/UwUshkin/task-3/internal/config"
	"github.com/UwUshkin/task-3/internal/processor"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to the YAML configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		wrappedErr := fmt.Errorf("fatal error loading config file '%s': %w", configPath, err)
		fmt.Fprintf(os.Stderr, "%v\n", wrappedErr)
		os.Exit(1)
	}

	if err := processor.ProcessAndSave(cfg.InputFile, cfg.OutputFile); err != nil {
		wrappedErr := fmt.Errorf("fatal error during data processing: %w", err)
		fmt.Fprintf(os.Stderr, "%v\n", wrappedErr)
		os.Exit(1)
	}
}
