package xmldecoder

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tntkatz/task-3/internal/data"
	"golang.org/x/text/encoding/charmap"
)

var ErrUnsupportedCharset = errors.New("unsupported charset")

func DecodeXML(inputFile []byte, valCurs *data.ValCurs) error {
	decoder := xml.NewDecoder(bytes.NewReader(inputFile))

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, fmt.Errorf("failed to decode %w: %s", ErrUnsupportedCharset, charset)
	}

	err := decoder.Decode(valCurs)
	if err != nil {
		return fmt.Errorf("XML decode error: %w", err)
	}

	return nil
}
