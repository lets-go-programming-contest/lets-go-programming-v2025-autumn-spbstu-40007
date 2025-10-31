package data

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

const PermissionsMkdir = 0o755
const PermissionsOpenFile = 0o600

type Valutes struct {
	Valutes []Valute `xml:"Valute"`
}
type Valute struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    Float  `json:"value"     xml:"Value"`
}

type Float float32

func (float *Float) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	decodeElement := ""
	if err := d.DecodeElement(&decodeElement, &start); err != nil {
		return fmt.Errorf("Valutes: %w", err)
	}

	decodeElement = strings.ReplaceAll(decodeElement, ",", ".")

	f, err := strconv.ParseFloat(decodeElement, 32)
	if err != nil {
		return fmt.Errorf("Valutes: %w", err)
	}

	*float = Float(f)

	return nil
}

func New(path string) (*Valutes, error) {
	ValutesContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Valutes: %w", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(ValutesContent)))
	decoder.CharsetReader = charset.NewReaderLabel

	Valutes := &Valutes{
		Valutes: []Valute{},
	}
	if err = decoder.Decode(Valutes); err != nil {
		return nil, fmt.Errorf("Valutes: %w", err)
	}

	return Valutes, nil
}

func (Valutes *Valutes) SaveToOutputFile(path string) error {
	dirs := filepath.Dir(path)

	err := os.MkdirAll(dirs, PermissionsMkdir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Valutes: %w", err)
	}

	output, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, PermissionsOpenFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Valutes: %w", err)
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(Valutes.Valutes); err != nil {
		fmt.Fprintln(os.Stderr, "Valutes: %w", err)
	}

	return nil
}
