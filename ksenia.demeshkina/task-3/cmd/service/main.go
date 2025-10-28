package main

import (
	"flag"
	"log"
	"sort"

	"github.com/ksuah/task-3/internal/config"
	"github.com/ksuah/task-3/internal/writer"
	"github.com/ksuah/task-3/internal/xmlparser"
)

func main() {
	configPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		log.Panicf("Missing --config flag with path to YAML config")
	}

	cfg := config.LoadConfig(*configPath)

	valutes := xmlparser.LoadXML(cfg.InputFile)

	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})

	writer.SaveToJSON(cfg.OutputFile, valutes)
}
