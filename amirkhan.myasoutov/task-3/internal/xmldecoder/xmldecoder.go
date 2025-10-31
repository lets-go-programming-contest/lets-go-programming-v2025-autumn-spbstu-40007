package xmldecoder

import (
	"encoding/xml"
	"os"

	"github.com/ami0-0/task-3/internal/data"
)

func DecodeXML(filePath string) ([]data.Valute, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var valCurs data.ValCurs
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, err
	}

	return valCurs.Currencies, nil
}
