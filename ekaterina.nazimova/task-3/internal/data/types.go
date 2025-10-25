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
	XMLName  xml.Name `xml:"Valute"`
	ID       string   `xml:"ID,attr"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  int      `xml:"Nominal"`
	Name     string   `xml:"Name"`
	ValueStr string   `xml:"Value"`
}

type OutputCurrency struct {
	NumCode  string  `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func (v Valute) ConvertToOutput() OutputCurrency {
	valueCleaned := strings.Replace(v.ValueStr, ",", ".", 1)
	value, _ := strconv.ParseFloat(valueCleaned, 64)

	normalizedValue := value / float64(v.Nominal)

	return OutputCurrency{
		NumCode:  v.NumCode,
		CharCode: v.CharCode,
		Value:    normalizedValue,
	}
}
