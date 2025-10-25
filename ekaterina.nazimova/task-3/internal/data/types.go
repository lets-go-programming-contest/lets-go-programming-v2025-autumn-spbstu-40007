package data

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName xml.Name `xml:"Valute"`
	ID      string   `xml:"ID,attr"`

	CharCode string `xml:"CharCode"`
	NumCode  string `xml:"NumCode"`

	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	ValueStr string `xml:"Value"`
}

type OutputCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func (v Valute) ConvertToOutput() OutputCurrency {
	valueCleaned := strings.Replace(v.ValueStr, ",", ".", 1)
	value, _ := strconv.ParseFloat(valueCleaned, 64)

	normalizedValue := value / float64(v.Nominal)

	numCodeInt, _ := strconv.Atoi(v.NumCode)

	return OutputCurrency{
		NumCode:  numCodeInt,
		CharCode: v.CharCode,
		Value:    normalizedValue,
	}
}
