package processor

import (
	"fmt"

	"github.com/ami0-0/task-3/internal/data"
	"github.com/ami0-0/task-3/internal/vp"
	"github.com/ami0-0/task-3/internal/xmldecoder"
)

func DecodeXMLFile(filePath string) ([]data.Valute, error) {
	currencies, err := xmldecoder.DecodeXML(filePath)
	if err != nil {
		return nil, fmt.Errorf("decode xml file: %w", err)
	}

	return currencies, nil
}

func SortCurrenciesByValue(currencies []data.Valute) []data.CurrencyOutput {
	return vp.SortAndConvert(currencies)
}

func SaveCurrenciesToJSON(currencies []data.CurrencyOutput, outputPath string) error {
	err := vp.SaveToJSON(currencies, outputPath)
	if err != nil {
		return fmt.Errorf("save to json: %w", err)
	}

	return nil
}
