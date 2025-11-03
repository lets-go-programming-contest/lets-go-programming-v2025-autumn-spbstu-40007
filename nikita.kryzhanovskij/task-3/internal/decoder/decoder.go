package decoder

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"nikita.kryzhanovskij/task-3/internal/models"
)

func DecodeXML(path string) (*models.ValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	defer func() {
		if errClose := file.Close(); errClose != nil {
			fmt.Fprintf(os.Stderr, "close file error: %v\n", errClose)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "windows-1251":
			return transform.NewReader(input, charmap.Windows1251.NewDecoder()), nil
		default:
			return input, nil
		}
	}

	var valCurs models.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &valCurs, nil
}
