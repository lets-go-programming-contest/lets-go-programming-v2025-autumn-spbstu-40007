package processor

import (
	"github.com/ami0-0/task-3/internal/data"
	"github.com/ami0-0/task-3/internal/vp"
	"github.com/ami0-0/task-3/internal/xmldecoder"
)

func DecodeXMLFile(filePath string) ([]data.Valute, error) {

	return xmldecoder.DecodeXML(filePath)
}

func SortCurrenciesByValue(currencies []data.Valute) []data.CurrencyOutput {
	
	return vp.SortAndConvert(currencies)
}

func SaveCurrenciesToJSON(currencies []data.CurrencyOutput, outputPath string) error {
	
	return vp.SaveToJSON(currencies, outputPath)
}
