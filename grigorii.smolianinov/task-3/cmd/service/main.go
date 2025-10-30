package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Smolyaninoff/GoLang/internal/config"
	"github.com/Smolyaninoff/GoLang/internal/currencies"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("ERROR: ")

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "Error: config path is required")
		os.Exit(1)
	}

	cfg := config.LoadConfig(*configPath)

	valutes := currencies.LoadAndSort(cfg.InputFile)

	currencies.SaveToJSON(cfg.OutputFile, valutes)

	fmt.Println("Processing completed successfully.")
}
