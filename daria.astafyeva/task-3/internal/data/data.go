package data

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

var ErrUnknownCharset = errors.New("unknown charset")

type Currency struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	ValueStr string   `xml:"Value"`
	Rate     float64  `xml:"-"`
}

type OutputCurrency struct {
	Num    int     `json:"num_code" xml:"NumCode" yaml:"num_code"`
	Char   string  `json:"char_code" xml:"CharCode" yaml:"char_code"`
	Amount float64 `json:"value" xml:"Value" yaml:"value"`
}

func LoadCurrencies(path string) []Currency {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("cannot read XML file '%s': %w", path, err))
	}

	decoder := xml.NewDecoder(strings.NewReader(string(content)))
	decoder.CharsetReader = func(charset string, r io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(r), nil
		}
		return nil, fmt.Errorf("%w: %s", ErrUnknownCharset, charset)
	}

	var wrapper struct {
		Currencies []Currency `xml:"Valute"`
	}
	if err := decoder.Decode(&wrapper); err != nil {
		panic(fmt.Errorf("XML decode failed: %w", err))
	}

	var list []Currency
	for _, c := range wrapper.Currencies {
		c.ValueStr = strings.ReplaceAll(c.ValueStr, ",", ".")
		val, err := strconv.ParseFloat(c.ValueStr, 64)
		if err != nil {
			panic(fmt.Errorf("invalid value '%s': %w", c.ValueStr, err))
		}
		c.Rate = val
		list = append(list, c)
	}
	return list
}
