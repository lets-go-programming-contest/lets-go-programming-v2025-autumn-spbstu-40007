package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Seeany1/task-3/internal/config"
	"github.com/Seeany1/task-3/internal/currencies"
)

func main() {
	cfgPath := "config.yaml"

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		cfgPath = "../../config.yaml"
	}

	cfg, err := config.LoadConfig(cfgPath)

	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	data, err := os.ReadFile(cfg.InputFile)

	if err != nil {
		log.Fatalf("failed to read xml %v", err)
	}

	svc := currencies.NewCurrencyService()
	list, err := svc.ParseXML(data)

	if err != nil {
		log.Fatalf("failed to parse xml %v", err)
	}

	svc.SortByValue(list)

	if err := svc.SaveToJSON(cfg.OutputFile, list); err != nil {
		log.Fatalf("failed to save json %v", err)
	}

	fmt.Printf("done output %s\n", cfg.OutputFile)
}
