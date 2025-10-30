package main

import (
	"flag"
	"log"
	"sort"

	"currency-processor/config"
	"currency-processor/data"
	"currency-processor/output"
)

func execute() {
	var settingsPath, resultFormat string

	flag.StringVar(&settingsPath, "config", "config.yaml", "YAML configuration file path")
	flag.StringVar(&resultFormat, "output-format", "json", "Result format (json/yaml/xml)")

	flag.Parse()

	appConfig := config.ReadSettings(settingsPath)

	currencyItems := data.ParseCurrencyFile(appConfig.SourcePath)

	sort.Sort(data.ByValueDesc(currencyItems))

	finalData := transformData(currencyItems)

	output.WriteProcessedData(appConfig.TargetPath, resultFormat, finalData)

	log.Printf(
		"Processing completed: %d currencies saved to '%s' in %s format",
		len(finalData), appConfig.TargetPath, resultFormat,
	)
}

func transformData(items []data.Currency) []data.ProcessedCurrency {
	transformed := make([]data.ProcessedCurrency, len(items))
	for index, item := range items {
		transformed[index] = item.ConvertToProcessed()
	}
	return transformed
}

func main() {
	execute()
}
