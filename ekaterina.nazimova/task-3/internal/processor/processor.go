package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/UwUshkin/task-3/internal/data"
	"github.com/UwUshkin/task-3/internal/xmldecoder"
)

func ProcessAndSave(inputPath, outputPath string) error {
	valCurs, err := xmldecoder.DecodeCBRXML(inputPath)
	if err != nil {
		return fmt.Errorf("decoding XML from %q: %w", inputPath, err)
	}

	var outputCurrencies []data.OutputCurrency
	for _, valute := range valCurs.Valutes {
		outputCurrencies = append(outputCurrencies, valute.ConvertToOutput())
	}

	sort.Slice(outputCurrencies, func(i, j int) bool {
		return outputCurrencies[i].Value > outputCurrencies[j].Value
	})

	jsonData, err := json.MarshalIndent(outputCurrencies, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling results to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0o644); err != nil {
		return fmt.Errorf("writing output file %q: %w", outputPath, err)
	}

	return nil
}
