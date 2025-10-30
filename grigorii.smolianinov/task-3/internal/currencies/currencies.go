package currencies

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Valute struct {
	NumCode  string  `json:"num_code" xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value" xml:"Value"`
}

func (v *Valute) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias Valute
	aux := struct {
		*Alias
		ValueStr string `xml:"Value"`
	}{
		Alias: (*Alias)(v),
	}

	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}

	normalizedStr := strings.Replace(aux.ValueStr, ",", ".", 1)
	f, err := strconv.ParseFloat(normalizedStr, 64)

	if err != nil {
		log.Panicf("Invalid number format in XML: %s (%v)", aux.ValueStr, err)
	}

	v.Value = f

	return nil
}

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

func LoadAndSort(path string) []Valute {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("Cannot open XML file: %v", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unknown charset: %s", charset) //nolint:err113
	}

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		log.Panicf("Cannot parse XML: %v", err)
	}

	sort.Slice(valCurs.Valutes, func(i, j int) bool {
		return valCurs.Valutes[i].Value > valCurs.Valutes[j].Value
	})

	return valCurs.Valutes
}

func SaveToJSON(path string, valutes []Valute) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Panicf("Failed to create directory: %v", err)
	}

	jsonData, err := json.MarshalIndent(valutes, "", "  ")
	if err != nil {
		log.Panicf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(path, jsonData, 0600); err != nil {
		log.Panicf("Failed to write JSON: %v", err)
	}
}
