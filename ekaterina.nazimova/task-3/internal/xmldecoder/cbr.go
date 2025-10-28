package xmldecoder

import (
	"encoding/xml"
	"fmt"
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
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}

	var result data.ValCurs
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding XML structure: %w", err)
	}

	return &result, nil
}

func DecodeCBRXML(filePath string) (*data.ValCurs, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening XML file %q: %w", filePath, err)
	}

	defer func() {
		if closeErr := xmlFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return decodeXMLFromReader(xmlFile)
}
