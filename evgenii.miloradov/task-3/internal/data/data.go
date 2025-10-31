package data

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
	DirectoryPermissions = 0o755
	FilePermissions      = 0o600
)

type Amount float32

func (amount *Amount) UnmarshalXML(decoder *xml.Decoder, startElement xml.StartElement) error {
	var rawValue string
	if decodeErr := decoder.DecodeElement(&rawValue, &startElement); decodeErr != nil {
		return fmt.Errorf("currency amount decode failed: %w", decodeErr)
	}

	normalizedValue := strings.ReplaceAll(rawValue, ",", ".")
	parsedValue, parseErr := strconv.ParseFloat(normalizedValue, 32)
	if parseErr != nil {
		return fmt.Errorf("currency amount parse failed: %w", parseErr)
	}

	*amount = Amount(parsedValue)
	return nil
}

type Currency struct {
	NumericCode int    `json:"numeric_code" xml:"NumCode"`
	Symbol      string `json:"symbol"       xml:"CharCode"`
	Value       Amount `json:"value"        xml:"Value"`
}

type CurrencyCollection struct {
	Currencies []Currency `xml:"Valute"`
}

func LoadFromFile(filePath string) (*CurrencyCollection, error) {
	fileContent, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read currency file: %w", readErr)
	}

	xmlDecoder := xml.NewDecoder(strings.NewReader(string(fileContent)))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	collection := &CurrencyCollection{}

	if decodeErr := xmlDecoder.Decode(collection); decodeErr != nil {
		return nil, fmt.Errorf("failed to parse currency data: %w", decodeErr)
	}

	return collection, nil
}

func (collection *CurrencyCollection) ExportToFile(outputPath string) error {
	outputDir := filepath.Dir(outputPath)

	if mkdirErr := os.MkdirAll(outputDir, DirectoryPermissions); mkdirErr != nil {
		return fmt.Errorf("failed to create output directory: %w", mkdirErr)
	}

	outputFile, fileErr := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, FilePermissions)
	if fileErr != nil {
		return fmt.Errorf("failed to create output file: %w", fileErr)
	}
	defer outputFile.Close()

	jsonEncoder := json.NewEncoder(outputFile)
	jsonEncoder.SetIndent("", "  ")

	if encodeErr := jsonEncoder.Encode(collection.Currencies); encodeErr != nil {
		return fmt.Errorf("failed to encode currency data: %w", encodeErr)
	}

	return nil
}
