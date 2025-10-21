package xmldecoder

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/treadwave/task-3/internal/structs"
	"golang.org/x/net/html/charset"
)

func XmlDecode(filepath string) (structs.TempValutes, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return structs.TempValutes{}, fmt.Errorf("failed to open config file: %w", err)
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	tempValutes := structs.TempValutes{}

	err = decoder.Decode(&tempValutes)

	if err != nil {
		return structs.TempValutes{}, fmt.Errorf("failed to decode XML:  %w", err)
	}

	return tempValutes, nil

}
