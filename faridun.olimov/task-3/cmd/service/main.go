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

	cfg := config.LoadConfig(configPath)

	valutes := data.DecodeXMLData(cfg.InputFile)

	sort.Sort(data.CustomSorter(valutes))

	resultValutes := prepareResults(valutes)

	output.SaveResults(cfg.OutputFile, outputFormat, resultValutes)

	fmt.Printf("Successfully saved %d currencies to file '%s' in '%s' format.\n", len(resultValutes), cfg.OutputFile, outputFormat)
}

func prepareResults(valutes []data.Valute) []data.ResultValute {
	resultValutes := make([]data.ResultValute, len(valutes))
	for i, v := range valutes {
		resultValutes[i] = v.ToResultValute()
	}
	
	return resultValutes
}
