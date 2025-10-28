package data

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (c *CurrencyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var xmlString string
	if err := d.DecodeElement(&xmlString, &start); err != nil {
		return fmt.Errorf("decoding element: %w", err)
	}

	cleanString := strings.Replace(xmlString, ",", ".", 1)
	parsedFloat, err := strconv.ParseFloat(cleanString, 64)
	if err != nil {
		return fmt.Errorf("parsing cleaned float %q: %w", cleanString, err)
	}

	*c = CurrencyValue(parsedFloat)
	return nil
}

type Valute struct {
	ID         string `xml:"ID,attr"    json:"-"`
	NominalStr string `xml:"Nominal"    json:"-"`
	Name       string `xml:"Name"       json:"-"`

	CharCode string `json:"char_code" xml:"CharCode"`
	NumCode  int    `json:"num_code"  xml:"NumCode"`

	Value CurrencyValue `json:"value"     xml:"Value"`
}

type ValCurs struct {
	Date    string       `xml:"Date,attr"`
	Name    string       `xml:"name,attr"`
	Valutes CurrencyList `xml:"Valute"`
}

type CurrencyList []Valute

func (c CurrencyList) Len() int {
	return len(c)
}

func (c CurrencyList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CurrencyList) Less(i, j int) bool {
	return c[i].Value > c[j].Value
}
