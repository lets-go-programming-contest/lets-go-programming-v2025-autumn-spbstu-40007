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
		return fmt.Errorf("failed to decode element: %w", err)
	}

	newValue := strings.Replace(xmlCode, ",", ".", 1)

	parsedValue, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return fmt.Errorf("failed to parse float: %w", err)
	}

	*f = FloatValue(parsedValue)

	return nil
}

type Valute struct {
	NumCode  int        `xml:"NumCode"  json:"num_code"`
	CharCode string     `xml:"CharCode" json:"char_code"`
	Value    FloatValue `xml:"Value"    json:"value"`
}

type ValCurs struct {
	Date   string   `xml:"Date,attr"`
	Name   string   `xml:"name,attr"`
	Valute []Valute `xml:"Valute"`
}

type ByValue []Valute

func (a ByValue) Len() int {
	return len(a)
}

func (a ByValue) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByValue) Less(i, j int) bool {
	return a[i].Value > a[j].Value
}
