package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"task-3/internal/config"
	"task-3/internal/currencies"
)

func main() {
	configPath := flag.String("config", "", "Path to config")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "Error: config path is required")
		os.Exit(1)
	}

	config, err := config.New(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to load config: %w", err)
		os.Exit(1)
	}

	currencies, err := currencies.New(config.InputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to load currencies: %w", err)
		os.Exit(1)
	}

	sort.Slice(currencies.Currencies, func(i, j int) bool {
		return currencies.Currencies[i].Value > currencies.Currencies[j].Value
	})

	if err := currencies.SaveToOutputFile(config.OutputFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to save to output file: %w", err)
		os.Exit(1)
	}
}
