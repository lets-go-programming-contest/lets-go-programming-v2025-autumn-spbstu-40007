package data

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type FloatValue float64

func (f *FloatValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var xmlCode string

	if err := d.DecodeElement(&xmlCode, &start); err != nil {
		return fmt.Errorf("decoding element: %w", err)
	}

	newValue := strings.Replace(xmlCode, ",", ".", 1)

	parsedValue, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return fmt.Errorf("parsing float: %w", err)
	}

	*f = FloatValue(parsedValue)

	return nil
}

type Valute struct {
	NumCode  int        `json:"num_code"  xml:"NumCode"`
	CharCode string     `json:"char_code" xml:"CharCode"`
	Value    FloatValue `json:"value"     xml:"Value"`
}

type ValCurs struct {
	Date    string  `xml:"Date,attr"`
	Name    string  `xml:"name,attr"`
	Valutes Valutes `xml:"Valute"`
}

type Valutes []Valute

func (v Valutes) Len() int {
	return len(v)
}

func (v Valutes) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Valutes) Less(i, j int) bool {
	return v[i].Value > v[j].Value
}
