package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/itsdasha/task-3/internal/config"
	"github.com/itsdasha/task-3/internal/data"
	"github.com/itsdasha/task-3/internal/output"
)

func main() {
	configPathPtr := flag.String("config", "", "Path to the YAML configuration file")
	outputFormatPtr := flag.String("output-format", "json", "Output file format (json, yaml, xml)")

	flag.Parse()

	if *configPathPtr == "" {
		log.Panicf("Missing --config flag with path to YAML config")
	}

	cfg := config.LoadConfig(*configPathPtr)

	valutes := data.DecodeXMLData(cfg.InputFile)

	sort.Sort(data.ByValueDesc(valutes))

	resultValutes := make([]data.ResultValute, len(valutes))
	for i, v := range valutes {
		resultValutes[i] = v.ToResultValute()
	}

	output.SaveResults(cfg.OutputFile, *outputFormatPtr, resultValutes)

	fmt.Printf(
		"Successfully processed %d currencies and saved to '%s' in '%s' format.\n",
		len(resultValutes), cfg.OutputFile, *outputFormatPtr,
	)
}
