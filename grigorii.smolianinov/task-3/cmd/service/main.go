package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/myuser/go-task/internal/config"
	"github.com/myuser/go-task/internal/writer"
	"github.com/myuser/go-task/internal/xmlparser"
)

func main() {
	configPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		log.Panic("Missing --config flag with path to YAML config")
	}

	cfg := config.LoadConfig(*configPath)

	valutes := xmlparser.LoadXML(cfg.InputFile)

	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})

	writer.SaveToJSON(cfg.OutputFile, valutes)

	fmt.Printf("Обработка завершена. Результат сохранен в %s. Количество валют: %d\n", cfg.OutputFile, len(valutes))
}
