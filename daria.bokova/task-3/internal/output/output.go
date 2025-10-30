package output

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bdshka/task-3/internal/data"
	"gopkg.in/yaml.v3"
)

var ErrorFormatNotSupported = errors.New("неподдерживаемый формат вывода")

type DataExporter struct{}

func ExportData(filename string, formatType string, currencies []data.ProcessedCurrency) {
	outputData, encodeErr := encodeToFormat(formatType, currencies)
	if encodeErr != nil {
		handleEncodingError(formatType, encodeErr)
	}

	dirErr := createOutputDirectory(filename)
	if dirErr != nil {
		handleDirectoryError(filepath.Dir(filename), dirErr)
	}

	writeErr := writeToFile(filename, outputData)
	if writeErr != nil {
		handleFileWriteError(filename, writeErr)
	}
}

func encodeToFormat(format string, currencies []data.ProcessedCurrency) ([]byte, error) {
	format = strings.ToLower(format)

	switch format {
	case "json":
		return encodeToJSON(currencies)
	case "yaml":
		return encodeToYAML(currencies)
	case "xml":
		return encodeToXML(currencies)
	default:
		return nil, fmt.Errorf("%w: %s. Доступные: json, yaml, xml",
			ErrorFormatNotSupported, format)
	}
}

func encodeToJSON(currencies []data.ProcessedCurrency) ([]byte, error) {
	jsonData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("ошибка кодирования JSON: %w", err)
	}
	return jsonData, nil
}

func encodeToYAML(currencies []data.ProcessedCurrency) ([]byte, error) {
	yamlData, err := yaml.Marshal(currencies)
	if err != nil {
		return nil, fmt.Errorf("ошибка кодирования YAML: %w", err)
	}
	return yamlData, nil
}

func encodeToXML(currencies []data.ProcessedCurrency) ([]byte, error) {
	xmlData := data.ProcessedCurrencyList{
		XMLName: xml.Name{Space: "", Local: "ValCurs"},
		Items:   currencies,
	}

	xmlBytes, err := xml.MarshalIndent(xmlData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("ошибка кодирования XML: %w", err)
	}

	return []byte(xml.Header + string(xmlBytes)), nil
}

func createOutputDirectory(filePath string) error {
	directory := filepath.Dir(filePath)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("невозможно создать директорию: %w", err)
	}
	return nil
}

func writeToFile(filePath string, data []byte) error {
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}
	return nil
}

func handleEncodingError(format string, err error) {
	fmt.Printf("Ошибка кодирования данных в формат '%s': %v\n", format, err)
	panic(fmt.Errorf("сбой преобразования данных: %w", err))
}

func handleDirectoryError(dirPath string, err error) {
	fmt.Printf("Ошибка создания директории '%s': %v\n", dirPath, err)
	panic(fmt.Errorf("невозможно создать директорию: %w", err))
}

func handleFileWriteError(filePath string, err error) {
	fmt.Printf("Ошибка записи в файл '%s': %v\n", filePath, err)
	panic(fmt.Errorf("сбой сохранения файла: %w", err))
}
