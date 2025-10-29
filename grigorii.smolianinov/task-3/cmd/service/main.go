package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Smolyaninoff/GoLang.git/task-3/internal/config"
	"github.com/Smolyaninoff/GoLang.git/task-3/internal/sorter"
	"github.com/Smolyaninoff/GoLang.git/task-3/internal/writer"
	"github.com/Smolyaninoff/GoLang.git/task-3/internal/xmlparser"
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
