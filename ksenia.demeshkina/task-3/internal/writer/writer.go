package writer

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/ksuah/task-3/internal/xmlparser"
)

type CurrencyOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func SaveToJSON(path string, valutes []xmlparser.Valute) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Panicf("Failed to create directory: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		log.Panicf("Cannot create output file: %v", err)
	}

	defer file.Close()

	var output []CurrencyOutput
	for _, v := range valutes {
		output = append(output, CurrencyOutput{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    v.Value,
		})
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		log.Panicf("Failed to write JSON: %v", err)
	}
}
