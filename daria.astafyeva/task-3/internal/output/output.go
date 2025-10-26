package output

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/itsdasha/task-3/internal/data"
	"gopkg.in/yaml.v3"
)

func SaveResults(filePath string, format string, result []data.ResultValute) {
	var encodedData []byte
	var marshalError error

	switch strings.ToLower(format) {
	case "json":
		encodedData, marshalError = json.MarshalIndent(result, "", "  ")
	case "yaml":
		encodedData, marshalError = yaml.Marshal(result)
	case "xml":

		resultXML := data.ResultValutes{Valutes: result}
		encodedData, marshalError = xml.MarshalIndent(resultXML, "", "  ")
		if marshalError == nil {
			encodedData = []byte(xml.Header + string(encodedData))
		}
	default:
		panic(fmt.Errorf("unsupported output format: %s. Available: json, yaml, xml", format))
	}

	if marshalError != nil {
		panic(fmt.Errorf("data encoding error for format '%s': %w", format, marshalError))
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(fmt.Errorf("failed to create directory '%s': %w", dir, err))
	}

	if err := os.WriteFile(filePath, encodedData, 0o600); err != nil {
		panic(fmt.Errorf("failed to write result to file '%s': %w", filePath, err))
	}
}
