package main

import (
	"flag"
	"log"
	"os"

	"github.com/falsefeelings/task-3/iternal/config"
	"github.com/falsefeelings/task-3/iternal/converter"
	"github.com/falsefeelings/task-3/iternal/data"
	"github.com/falsefeelings/task-3/iternal/pathmaker"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to config file")
	flag.Parse()

	config, err := config.Read(configPath)
	if err != nil {
		log.Fatal("Config error:", err)
	}

	if _, err := os.Stat(config.InputFile); os.IsNotExist(err) {
		log.Fatalf("no such file or directorty")
	}

	xmlData, err := os.ReadFile(config.InputFile)
	if err != nil {
		log.Fatal("Read file error:", err)
	}

	valCourse, err := data.ParseXML(xmlData)
	if err != nil {
		log.Fatal("XML parse error:", err)
	}

	conv := converter.New()
	outputData, err := conv.Convert(valCourse, config.OutputFormat)
	if err != nil {
		log.Fatal("Conversion error:", err)
	}

	err = pathmaker.CreateOutPath(config.OutputFile)
	if err != nil {
		log.Fatal("Create path error:", err)
	}

	err = os.WriteFile(config.OutputFile, outputData, 0644)
	if err != nil {
		log.Fatal("Write file error:", err)
	}
}
