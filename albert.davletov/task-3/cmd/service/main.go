package main

import (
	"flag"
	"log"

	"github.com/treadwave/task-3/internal/config"
	"github.com/treadwave/task-3/internal/converter"
	"github.com/treadwave/task-3/internal/jsonencoder"
	"github.com/treadwave/task-3/internal/xmldecoder"
)

func main() {
	configPath := flag.String("config", "", "enter a path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		log.Panic("no config file")
	}

	configStruct, err := config.YamlDecoder(*configPath)
	if err != nil {
		log.Panic(err)
	}

	xmlValutes, err := xmldecoder.XmlDecode(configStruct.InputFile)
	if err != nil {
		log.Panic(err)
	}

	convertedValutes, err := converter.Converter(xmlValutes.TempValutes)
	if err != nil {
		log.Panic(err)
	}

	err = jsonencoder.JsonEncoder(convertedValutes, configStruct.OutputFile)
	if err != nil {
		log.Panic(err)
	}
}
