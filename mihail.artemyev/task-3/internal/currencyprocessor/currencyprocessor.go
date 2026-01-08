package currencyprocessor

import (
	"fmt"
	"os"
	"sort"

	"github.com/Mart22052006/task-3/internal/configuration"
	"github.com/Mart22052006/task-3/internal/currencydata"
	"github.com/Mart22052006/task-3/internal/outputmanager"
	"github.com/Mart22052006/task-3/internal/xmlparser"
)

func ExecuteProcessing(applicationConfig *configuration.ApplicationConfig) error {
	inputFileContent, fileReadError := os.ReadFile(applicationConfig.InputFilePath)
	if fileReadError != nil {
		return fmt.Errorf("reading input file from %q: %w", applicationConfig.InputFilePath, fileReadError)
	}

	exchangeRateData := &currencydata.ExchangeRateList{
		ExchangeDate: "",
		MarketName:   "",
		CurrencyList: []currencydata.CurrencyExchange{},
	}

	xmlParsingError := xmlparser.ParseXMLFromBytes(inputFileContent, exchangeRateData)
	if xmlParsingError != nil {
		return fmt.Errorf("parsing XML data: %w", xmlParsingError)
	}

	sortedCurrencyList := make(currencydata.CurrencyCollection, len(exchangeRateData.CurrencyList))
	copy(sortedCurrencyList, exchangeRateData.CurrencyList)

	sort.Sort(sortedCurrencyList)

	jsonOutputError := outputmanager.WriteJSONOutput(
		applicationConfig.OutputFilePath,
		[]currencydata.CurrencyExchange(sortedCurrencyList),
	)
	if jsonOutputError != nil {
		return fmt.Errorf("writing JSON output: %w", jsonOutputError)
	}

	return nil
}
