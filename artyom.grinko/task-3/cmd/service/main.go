package main

import (
	"cmp"
	"flag"
	"slices"
	"sort"

	"task-3/internal/config"
	"task-3/internal/die"
	"task-3/internal/exchangeRate"
)

func main() {
	configPath := flag.String("config", "", "Path to config")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		die.Die("Path to config was not provided")
	}

	config, err := config.FromFile(*configPath)
	if err != nil {
		die.Die(err)
	}

	rate, err := exchangeRate.FromXMLFile(config.InputFile)
	if err != nil {
		die.Die(err)
	}

	slices.SortFunc(rate.Currencies, func(x, y exchangeRate.Currency) int {
		return cmp.Compare(x.Value, y.Value)
	})

	if err = rate.ToJSONFile(config.OutputFile); err != nil {
		die.Die(err)
	}
}
