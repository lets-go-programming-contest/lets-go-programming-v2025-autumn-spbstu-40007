package main

//nolint:gofumpt
import (
	"flag"
	"fmt"
	"os"
	"sort"

	"task-3/internal/config"
	"task-3/internal/currencies"
)

func main() {
	confPath := flag.String("config", "", "Path to config")
	flag.Parse()

	if *confPath == "" {
		fmt.Fprintln(os.Stderr, "Config path is required")
		os.Exit(1)
	}

	cfg, err := config.LoadCurrencies(*confPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can not load config: did not find expected key %w", err)
		os.Exit(1)
	}

	currenciesData, err := currencies.New(cfg.InputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can not load currencies: %w", err)
		os.Exit(1)
	}

	sort.Slice(currenciesData.Currencies, func(i, j int) bool {
		return currenciesData.Currencies[i].Value > currenciesData.Currencies[j].Value
	})

	if err := currenciesData.WriteToFile(cfg.OutputFile); err != nil {
		fmt.Fprintln(os.Stderr, "Can not save output: %w", err)
		os.Exit(1)
	}
}
