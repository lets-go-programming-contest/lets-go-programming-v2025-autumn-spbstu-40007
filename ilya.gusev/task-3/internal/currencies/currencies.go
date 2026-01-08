package currencies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

const (
	dirPermissions  = 0o755
	filePermissions = 0o644
)

type ExchangeRate float64

func (rate *ExchangeRate) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var rawValue string
	if err := decoder.DecodeElement(&rawValue, &start); err != nil {
		return fmt.Errorf("failed to decode xml element: %w", err)
	}

	normalizedValue := strings.ReplaceAll(rawValue, ",", ".")

	parsedFloat, err := strconv.ParseFloat(normalizedValue, 64)
	if err != nil {
		return fmt.Errorf("failed to parse exchange rate: %w", err)
	}

	*rate = ExchangeRate(parsedFloat)

	return nil
}

type CurrencyItem struct {
	Code   int          `json:"num_code" xml:"NumCode"`
	Symbol string       `json:"char_code" xml:"CharCode"`
	Rate   ExchangeRate `json:"value" xml:"Value"`
}

type CurrencyData struct {
	Items []CurrencyItem `xml:"Valute"`
}

func LoadFromXML(filePath string) (*CurrencyData, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read xml file: %w", err)
	}

	xmlDecoder := xml.NewDecoder(strings.NewReader(string(fileContent)))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	var data CurrencyData

	if err := xmlDecoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("cannot decode xml data: %w", err)
	}

	return &data, nil
}

func SaveToJSON(filePath string, items []CurrencyItem) error {
	directory := filepath.Dir(filePath)

	if err := os.MkdirAll(directory, dirPermissions); err != nil {
		return fmt.Errorf("cannot create output directory: %w", err)
	}

	outputFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePermissions)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(items); err != nil {
		return fmt.Errorf("cannot encode json data: %w", err)
	}

	return nil
}
