package xmldecoder

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/ami0-0/task-3/internal/data"

	"golang.org/x/text/encoding/charmap"
)

func DecodeXML(filePath string) ([]data.Valute, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var valCurs data.ValCurs
	decoder := xml.NewDecoder(file)

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return input, nil
		}
	}

	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, err
	}

	return valCurs.Currencies, nil
}
