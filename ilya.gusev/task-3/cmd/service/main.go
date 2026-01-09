package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"task-3/internal/config"
	"task-3/internal/currencies"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("error: config path must be specified")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	data, err := currencies.LoadFromXML(cfg.InputFile)
	if err != nil {
		log.Fatalf("failed to load currency data: %v", err)
	}

	sort.Slice(data.Items, func(i, j int) bool {
		return data.Items[i].Rate > data.Items[j].Rate
	})

	if err := currencies.SaveToJSON(cfg.OutputFile, data.Items); err != nil {
		log.Fatalf("failed to save results: %v", err)
	}

	fmt.Println("Currency data processed successfully")
}
