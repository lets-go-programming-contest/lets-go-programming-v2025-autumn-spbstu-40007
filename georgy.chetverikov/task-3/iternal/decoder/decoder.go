package decoder

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/falsefeelings/task-3/iternal/data"
	"golang.org/x/text/encoding/charmap"
)

var ErrUnsupportedCharset = errors.New("unsupported charset")

func DecodeXML(inputFile []byte, valCurls *data.Valute) error {
	decoder := xml.NewDecoder(bytes.NewReader(inputFile))

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, fmt.Errorf("decoding %w: %q", ErrUnsupportedCharset, charset)
	}

	err := decoder.Decode(valCurls)
	if err != nil {
		return fmt.Errorf("decoding XML: %w", err)
	}

	return nil
}
