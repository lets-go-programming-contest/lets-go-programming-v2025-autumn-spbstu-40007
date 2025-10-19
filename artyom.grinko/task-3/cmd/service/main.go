package main

//nolint:gofumpt
import (
	"cmp"
	"flag"
	"slices"

	"task-3/internal/config"
	"task-3/internal/die"
	"task-3/internal/exchangerate"
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

	rate, err := exchangerate.FromXMLFile(config.InputFile)
	if err != nil {
		die.Die(err)
	}

	slices.SortFunc(rate.Currencies, func(x, y exchangerate.Currency) int {
		return cmp.Compare(y.Value, x.Value)
	})

	if err = rate.ToJSONFile(config.OutputFile); err != nil {
		die.Die(err)
	}
}
