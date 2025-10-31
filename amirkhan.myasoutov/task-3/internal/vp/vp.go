package vp

import (
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ami0-0/task-3/internal/data"
)

func SortAndConvert(currencies []data.Valute) []data.CurrencyOutput {
	output := make([]data.CurrencyOutput, len(currencies))

	for i, currency := range currencies {
		valueStr := strings.Replace(currency.Value, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			panic(err)
		}

		output[i] = data.CurrencyOutput{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    value,
		}
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Value > output[j].Value
	})

	return output
}

func SaveToJSON(currencies []data.CurrencyOutput, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(currencies)
}
