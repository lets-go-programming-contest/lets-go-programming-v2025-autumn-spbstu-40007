package output

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/Maska192/task-3/internal/data"
)

func saveResults(filePath string, format string, data []data.ResultValute) {
	var encodedData []byte
	var err error

	switch strings.ToLower(format) {
	case "json":
		encodedData, err = json.MarshalIndent(data, "", "  ")
	case "yaml":
		encodedData, err = yaml.Marshal(data)
	case "xml":
		resultXML := ResultValutes{Valutes: data}
		encodedData, err = xml.MarshalIndent(resultXML, "", "  ")
		encodedData = []byte(xml.Header + string(encodedData))
	default:
		panic(fmt.Errorf("unsupported output format: %s. Available: json, yaml, xml", format))
	}

	if err != nil {
		fmt.Printf("Error encoding data to '%s' format: %v\n", format, err)
		panic(fmt.Errorf("data encoding error: %w", err))
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory '%s': %v\n", dir, err)
		panic(fmt.Errorf("failed to create directory: %w", err))
	}

	if err := os.WriteFile(filePath, encodedData, 0644); err != nil {
		fmt.Printf("Error writing to file '%s': %v\n", filePath, err)
		panic(fmt.Errorf("failed to write result to file: %w", err))
	}
}
