package main

import (
	"flag"

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

	exchangeRate, err := exchangeRate.FromXMLFile(config.InputFile)
	if err != nil {
		die.Die(err)
	}

	if err = exchangeRate.ToJSONFile(config.OutputFile); err != nil {
		die.Die(err)
	}
}
