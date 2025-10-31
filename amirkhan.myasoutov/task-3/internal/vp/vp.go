package vp

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ami0-0/task-3/internal/data"
)

func SortAndConvert(currencies []data.Valute) []data.CurrencyOutput {
	output := make([]data.CurrencyOutput, len(currencies))

	for index, currency := range currencies {
		valueStr := strings.ReplaceAll(currency.Value, ",", ".")
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			panic(err)
		}

		output[index] = data.CurrencyOutput{
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
		return fmt.Errorf("create file: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(currencies)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("encode json: %w", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("close file: %w", err)
	}

	return nil
}
