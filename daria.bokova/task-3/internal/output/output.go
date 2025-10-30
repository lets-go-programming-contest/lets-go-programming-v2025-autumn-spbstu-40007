package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"currency-converter/internal/data"
)

func SaveCurrencies(currencies []data.Currency, outputFile string) error {
	dir := filepath.Dir(outputFile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания директории: %v", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(currencies)
	if err != nil {
		return fmt.Errorf("ошибка кодирования JSON: %v", err)
	}

	return nil
}
