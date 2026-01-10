package xmlparser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Mart22052006/task-3/internal/currencydata"
	"golang.org/x/text/encoding/charmap"
)

var ErrUnsupportedCharsetError = errors.New("unsupported character set encoding")

func ParseXMLFromBytes(
	xmlFileContent []byte,
	exchangeRateListTarget *currencydata.ExchangeRateList,
) error {
	xmlDecoderInstance := xml.NewDecoder(bytes.NewReader(xmlFileContent))

	xmlDecoderInstance.CharsetReader = func(charsetName string, charsetInputStream io.Reader) (io.Reader, error) {
		if strings.ToLower(charsetName) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(charsetInputStream), nil
		}

		return nil, fmt.Errorf("charset %q is not supported: %w", charsetName, ErrUnsupportedCharsetError)
	}

	xmlDecodingError := xmlDecoderInstance.Decode(exchangeRateListTarget)
	if xmlDecodingError != nil {
		return fmt.Errorf("decoding XML content: %w", xmlDecodingError)
	}

	return nil
}
