package writer

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/Smolyaninoff/GoLang.git/task-3/internal/xmlparser"
)

type CurrencyOutput struct {
	NumCode  string  `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func SaveToJSON(path string, valutes []xmlparser.Valute) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Panicf("Failed to create directory: %v", err)
	}

	output := make([]CurrencyOutput, 0, len(valutes))
	for _, v := range valutes {
		output = append(output, CurrencyOutput{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    v.Value,
		})
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Panicf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		log.Panicf("Failed to write JSON: %v", err)
	}
}
