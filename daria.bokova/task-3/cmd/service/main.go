package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/bdshka/task-3/internal/config"
	"github.com/bdshka/task-3/internal/data"
	"github.com/bdshka/task-3/internal/output"
)

type AppConfig struct {
	ConfigFile string
	OutputType string
}

type CurrencyProcessor struct{}

type sortAdapter struct {
	data.CurrencySorter
}

func (a sortAdapter) Len() int           { return a.Count() }
func (a sortAdapter) Swap(i, j int)      { a.Exchange(i, j) }
func (a sortAdapter) Less(i, j int) bool { return a.Compare(i, j) }

func main() {
	appConfig := setupFlags()

	configuration := config.LoadSettings(appConfig.ConfigFile)

	currencyList := data.ProcessXMLFile(configuration.InputFile)

	sort.Sort(sortAdapter{currencyList})

	processedCurrencies := processCurrencyData(currencyList)

	output.ExportData(
		configuration.OutputFile,
		appConfig.OutputType,
		processedCurrencies,
	)

	displaySuccessMessage(len(processedCurrencies), configuration.OutputFile, appConfig.OutputType)
}

func setupFlags() *AppConfig {
	configFilePath := flag.String("config", "config.yaml", "Путь к YAML файлу конфигурации")
	outputType := flag.String("output-format", "json", "Формат выходного файла (json, yaml, xml)")

	flag.Parse()

	return &AppConfig{
		ConfigFile: *configFilePath,
		OutputType: *outputType,
	}
}

func processCurrencyData(currencies []data.CurrencyItem) []data.ProcessedCurrency {
	results := make([]data.ProcessedCurrency, len(currencies))

	for index, currency := range currencies {
		results[index] = currency.ConvertToOutputFormat()
	}

	return results
}

func displaySuccessMessage(count int, filename, format string) {
	fmt.Printf(
		"Успешно сохранено %d валют в файл '%s' в формате '%s'.\n",
		count, filename, format,
	)
}
