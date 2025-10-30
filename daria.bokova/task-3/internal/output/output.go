package output

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"strings"

	"currency-processor/data"

	"gopkg.in/yaml.v3"
)

func WriteProcessedData(filepath string, format string, currencies []data.ProcessedCurrency) {
	var serializedData []byte
	var serializationErr error

	switch strings.ToLower(format) {
	case "json":
		serializedData, serializationErr = json.MarshalIndent(currencies, "", "  ")
	case "yaml":
		serializedData, serializationErr = yaml.Marshal(currencies)
	case "xml":
		xmlData := data.CurrencyCollection{
			XMLName:    xml.Name{Space: "", Local: "ValCurs"},
			Currencies: currencies,
		}
		serializedData, serializationErr = xml.MarshalIndent(xmlData, "", "  ")
		serializedData = []byte(xml.Header + string(serializedData))
	default:
		panic("unsupported output format: " + format + ". Supported: json, yaml, xml")
	}

	if serializationErr != nil {
		log.Printf("Data serialization error to '%s': %v", format, serializationErr)
		panic("data serialization failed: " + serializationErr.Error())
	}

	targetDir := filepath.Dir(filepath)
	if dirErr := os.MkdirAll(targetDir, 0755); dirErr != nil {
		log.Printf("Directory creation error '%s': %v", targetDir, dirErr)
		panic("directory creation failed: " + dirErr.Error())
	}

	if writeErr := os.WriteFile(filepath, serializedData, 0600); writeErr != nil {
		log.Printf("File writing error '%s': %v", filepath, writeErr)
		panic("file writing failed: " + writeErr.Error())
	}
}
