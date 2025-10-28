package main

import (
	"flag"
	"fmt"
	"os"

	"nikita.kryzhanovskij/task-3/internal/config"
	"nikita.kryzhanovskij/task-3/internal/decoder"
	"nikita.kryzhanovskij/task-3/internal/encoder"
	"nikita.kryzhanovskij/task-3/internal/processor"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("Usage: program --config <path_to_config.yaml>")
		os.Exit(1)
	}

	if err := run(*configPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing completed successfully")
}

func run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	valCurs, err := decoder.DecodeXML(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	results, err := processor.Process(valCurs)
	if err != nil {
		return fmt.Errorf("process: %w", err)
	}

	if err := encoder.EncodeJSON(cfg.OutputFile, results); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
