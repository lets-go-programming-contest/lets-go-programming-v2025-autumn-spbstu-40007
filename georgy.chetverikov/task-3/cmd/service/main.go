package main

import (
	"encoding/xml"
	"log"
	"os"

	"github.com/falsefeelings/task-3/iternal/config"
	"github.com/falsefeelings/task-3/iternal/converter"
	"github.com/falsefeelings/task-3/iternal/data"
	"github.com/falsefeelings/task-3/iternal/pathmaker"
)

func main() {
	cfg, err := config.Read("config.yaml")
	if err != nil {
		log.Fatal("Config error:", err)
	}

	xmlData, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		log.Fatal("Read file error:", err)
	}

	var valCourse data.ValCourse
	err = xml.Unmarshal(xmlData, &valCourse)
	if err != nil {
		log.Fatal("XML parse error:", err)
	}

	conv := converter.New()
	outputData, err := conv.Convert(&valCourse, cfg.OutputFormat)
	if err != nil {
		log.Fatal("Conversion error:", err)
	}

	err = pathmaker.CreateOutPath(cfg.OutputFile)
	if err != nil {
		log.Fatal("Create path error:", err)
	}

	err = os.WriteFile(cfg.OutputFile, outputData, 0644)
	if err != nil {
		log.Fatal("Write file error:", err)
	}

	log.Printf("Success! Converted %s → %s (%s format)",
		cfg.InputFile, cfg.OutputFile, cfg.OutputFormat)
}
