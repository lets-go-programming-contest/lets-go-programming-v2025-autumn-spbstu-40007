package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Smolyaninoff/GoLang.git/internal/config"
	"github.com/Smolyaninoff/GoLang.git/internal/sorter"
	"github.com/Smolyaninoff/GoLang.git/internal/writer"
	"github.com/Smolyaninoff/GoLang.git/internal/xmlparser"
)

func main() {
	configPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		log.Panic("Missing --config flag with path to YAML config")
	}

	cfg := config.LoadConfig(*configPath)

	valutes := xmlparser.LoadXML(cfg.InputFile)

	sortedValutes := sorter.SortByValueDesc(valutes)

	writer.SaveToJSON(cfg.OutputFile, sortedValutes)

	fmt.Printf("✅ Успешно обработано %d валют. Результат сохранен в %s\n", len(sortedValutes), cfg.OutputFile)
}
