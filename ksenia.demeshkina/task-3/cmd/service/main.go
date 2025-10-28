package main

import (
	"flag"
	"log"
	"sort"

	"myapp/config"
	"myapp/writer"
	"myapp/xmlparser"
)

func main() {
	configPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		log.Panicf("Missing --config flag with path to YAML config")
	}

	cfg := config.LoadConfig(*configPath) // загружаем конфиг

	valutes := xmlparser.LoadXML(cfg.InputFile) // загружаем XML

	sort.Slice(valutes, func(i, j int) bool { // сортируем валюты
		return valutes[i].Value > valutes[j].Value
	})

	writer.SaveToJSON(cfg.OutputFile, valutes)

	log.Println("Result saved to:", cfg.OutputFile)
}
