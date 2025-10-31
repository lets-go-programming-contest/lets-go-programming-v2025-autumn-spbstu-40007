package xmldecoder

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/treadwave/task-3/internal/structs"
	"golang.org/x/net/html/charset"
)

func XMLDecode(filepath string) (structs.Valutes, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return structs.Valutes{}, fmt.Errorf("failed to open config file: %w", err)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing file: %w", err)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	Valutes := structs.Valutes{
		Valutes: []structs.Valute{},
	}

	err = decoder.Decode(&Valutes)
	if err != nil {
		return structs.Valutes{}, fmt.Errorf("failed to decode XML:  %w", err)
	}

	return Valutes, err
}
