package models

import "encoding/xml"

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID        string `xml:"ID,attr" json:"id"`
	NumCode   int    `xml:"NumCode" json:"num_code"`
	CharCode  string `xml:"CharCode" json:"char_code"`
	Nominal   int    `xml:"Nominal" json:"nominal"`
	Name      string `xml:"Name" json:"name"`
	Value     string `xml:"Value" json:"value"`
	VunitRate string `xml:"VunitRate" json:"vunit_rate"`
}

type ValuteOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}
