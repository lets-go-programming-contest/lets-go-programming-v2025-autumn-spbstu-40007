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

const (
	permMkdir = 0o755
	permFile  = 0o600
)

type Float float32

func (float *Float) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valStr string
	if err := d.DecodeElement(&valStr, &start); err != nil {
		return fmt.Errorf("decode xml value: %w", err)
	}
	valStr = strings.ReplaceAll(valStr, ",", ".")

	f64, err := strconv.ParseFloat(valStr, 32)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}
	*float = Float(f64)

	return nil
}

type Currency struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    Float  `json:"value"     xml:"Value"`
}

type Currencies struct {
	Currencies []Currency `xml:"Valute"`
}

func New(path string) (*Currencies, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	decoder.CharsetReader = charset.NewReaderLabel

	var cur Currencies
	if err := decoder.Decode(&cur); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &cur, nil
}

func (c *Currencies) WriteToFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, permMkdir); err != nil {
		return fmt.Errorf("create dirs: %w", err)
	}

	outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, permFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "currencies: %w", err)
	}

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c.Currencies); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
