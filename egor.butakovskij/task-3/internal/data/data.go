package data

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type FloatValue float64

func (f *FloatValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	newValue := strings.Replace(v, ",", ".", 1)

	parsedValue, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return err
	}

	*f = FloatValue(parsedValue)

	return nil
}

type Valute struct {
	ID        string     `xml:"ID,attr"`
	NumCode   int        `xml:"NumCode"`
	CharCode  string     `xml:"CharCode"`
	Nominal   string     `xml:"Nominal"`
	Name      string     `xml:"Name"`
	Value     FloatValue `xml:"Value"`
	VunitRate string     `xml:"VunitRate"`
}

type ValCurs struct {
	Date   string   `xml:"Date,attr"`
	Name   string   `xml:"name,attr"`
	Valute []Valute `xml:"Valute"`
}

type CurrencyResult struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
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
