package xmldecoder

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/ami0-0/task-3/internal/data"
	"golang.org/x/text/encoding/charmap"
)

func DecodeXML(filePath string) ([]data.Valute, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open xml file: %w", err)
	}

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return input, nil
		}
	}

	var valCurs data.ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	if err = file.Close(); err != nil {
		return nil, fmt.Errorf("close file: %w", err)
	}

	return valCurs.Currencies, nil
}
