package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/Maska192/task-3/internal/config"
	"github.com/Maska192/task-3/internal/data"
	"github.com/Maska192/task-3/internal/output"
)

var (
	configPath   string
	outputFormat string
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "Path to the YAML configuration file")
	flag.StringVar(&outputFormat, "output-format", "json", "Output file format (json, yaml, xml)")
}

func main() {
	flag.Parse()

	cfg := loadConfig(configPath)

	valutes := decodeXMLData(cfg.InputFile)

	sort.Sort(CustomSorter(valutes))

	resultValutes := prepareResults(valutes)

	saveResults(cfg.OutputFile, outputFormat, resultValutes)

	fmt.Printf("Successfully saved %d currencies to file '%s' in '%s' format.\n", len(resultValutes), cfg.OutputFile, outputFormat)
}

func prepareResults(valutes []Valute) []ResultValute {
	resultValutes := make([]ResultValute, len(valutes))
	for i, v := range valutes {
		resultValutes[i] = v.toResultValute()
	}
	return resultValutes
}
