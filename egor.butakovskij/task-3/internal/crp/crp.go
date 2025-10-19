package crp

import (
	"github.com/tntkatz/task-3/internal/data"
)

func CurrencyProcess(processedValutes []data.Valute) []data.CurrencyResult {
	currencyResults := make([]data.CurrencyResult, 0, len(processedValutes))

	for _, val := range processedValutes {
		currencyResult := data.CurrencyResult{
			NumCode:  val.NumCode,
			CharCode: val.CharCode,
			Value:    float64(val.Value),
		}

		currencyResults = append(currencyResults, currencyResult)
	}

	return currencyResults
}
