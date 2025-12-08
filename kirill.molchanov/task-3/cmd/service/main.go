package main

//nolint:gofumpt
import (
	"flag"
	"fmt"
	"os"
	"sort"

	"task-3/internal/config"
	"task-3/internal/data"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing config path")
		os.Exit(1)
	}

	config, err := config.ReadConfig(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	Valutes, err := data.New(config.InputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to load Valutes: %w", err)
		os.Exit(1)
	}

	sort.Slice(Valutes.Valutes, func(i, j int) bool {
		return Valutes.Valutes[i].Value > Valutes.Valutes[j].Value
	})

	if err := Valutes.SaveToOutputFile(config.OutputFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error: failed to save to output file: %w", err)
		os.Exit(1)
	}
}
