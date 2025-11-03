package main

import (
	"flag"
	"fmt"

	"github.com/se1lzor/task-3/internal/config"
	"github.com/se1lzor/task-3/internal/data"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("Error: Config file path is required. Use: -config config.yaml")
	}

	cfg, err := config.New(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}

	currencies, err := data.LoadFromXML(cfg.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Error loading currencies: %v", err))
	}

	data.SortByValue(currencies)

	err = data.SaveToJSON(currencies, cfg.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("Error saving JSON: %v", err))
	}

	fmt.Printf("Success! Processed %d currencies, saved to: %s\n",
		len(currencies), cfg.OutputFile)
}
