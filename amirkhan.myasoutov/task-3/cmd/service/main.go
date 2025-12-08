package main

import (
	"flag"

	"github.com/ami0-0/task-3/internal/config"
	"github.com/ami0-0/task-3/internal/pathcreator"
	"github.com/ami0-0/task-3/internal/processor"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("config path is required")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	err = pathcreator.EnsureDirectoryExists(cfg.OutputFile)
	if err != nil {
		panic(err)
	}

	currencies, err := processor.DecodeXMLFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sortedCurrencies := processor.SortCurrenciesByValue(currencies)

	err = processor.SaveCurrenciesToJSON(sortedCurrencies, cfg.OutputFile)
	if err != nil {
		panic(err)
	}
}
