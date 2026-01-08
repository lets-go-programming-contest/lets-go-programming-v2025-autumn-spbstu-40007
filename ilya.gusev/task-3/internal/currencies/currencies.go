package currencies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html/charset"
)

type ExchangeRate float64

func (e *ExchangeRate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valueStr string
	if err := d.DecodeElement(&valueStr, &start); err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	valueStr = strings.ReplaceAll(valueStr, ",", ".")
	var value float64
	if _, err := fmt.Sscanf(valueStr, "%f", &value); err != nil {
		return fmt.Errorf("failed to parse exchange rate: %w", err)
	}

	*e = ExchangeRate(value)
	return nil
}

func (e ExchangeRate) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(float64(e))
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}

	return data, nil
}

type CurrencyData struct {
	XMLName  xml.Name `xml:"Valute"`
	ID       string   `xml:"ID,attr"`
	Code     int      `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  int      `xml:"Nominal"`
	Name     string   `xml:"Name"`
	Rate     ExchangeRate
}

type CurrencyItem struct {
	Code     int          `json:"num_code" xml:"NumCode"`
	CharCode string       `json:"char_code" xml:"CharCode"`
	Rate     ExchangeRate `json:"value" xml:"Value"`
}

type ValCurs struct {
	XMLName xml.Name       `xml:"ValCurs"`
	Items   []CurrencyItem `xml:"Valute"`
}

func LoadFromXML(filePath string) (*ValCurs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data ValCurs
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return &data, nil
}

func SaveToJSON(filePath string, items []CurrencyItem) error {
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "failed to close output file: %v\n", closeErr)
		}
	}()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(items); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
