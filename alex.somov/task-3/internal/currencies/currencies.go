package currencies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Float float32

func (float *Float) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s := ""
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	s = strings.ReplaceAll(s, ",", ".")

	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	*float = Float(f)

	return nil
}

type Currency struct {
	NumCode  int    `xml:"NumCode" json:"num_code"`
	CharCode string `xml:"CharCode" json:"char_code"`
	Value    Float  `xml:"Value" json:"value"`
}

type Currencies struct {
	Currencies []Currency `xml:"Valute"`
}

func New(path string) (*Currencies, error) {
	currenciesContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("currencies: %w", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(currenciesContent)))
	decoder.CharsetReader = charset.NewReaderLabel

	currencies := &Currencies{}
	if err = decoder.Decode(currencies); err != nil {
		return nil, fmt.Errorf("currencies: %w", err)
	}

	return currencies, nil
}

func (currencies *Currencies) SaveToOutputFile(path string) error {
	dirs := filepath.Dir(path)
	err := os.MkdirAll(dirs, 0o755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "currencies: %w", err)
	}

	output, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "currencies: %w", err)
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(currencies.Currencies); err != nil {
		fmt.Fprintln(os.Stderr, "currencies: %w", err)
	}

	return nil
}
