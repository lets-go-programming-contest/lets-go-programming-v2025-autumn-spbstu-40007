package xmldecoder

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/UwUshkin/task-3/internal/data"
	"golang.org/x/text/encoding/charmap"
)

func decodeXMLFromReader(reader io.Reader) (*data.ValCurs, error) {
	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, nil
	}

	var valCurs data.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, err
	}

	return &valCurs, nil
}

func DecodeCBRXML(filePath string) (*data.ValCurs, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	return decodeXMLFromReader(xmlFile)
}
