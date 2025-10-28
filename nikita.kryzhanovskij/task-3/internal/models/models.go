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
	ID        string `json:"id"         xml:"ID,attr"`
	NumCode   int    `json:"num_code"   xml:"NumCode"`
	CharCode  string `json:"char_code"  xml:"CharCode"`
	Nominal   int    `json:"nominal"    xml:"Nominal"`
	Name      string `json:"name"       xml:"Name"`
	Value     string `json:"value"      xml:"Value"`
	VunitRate string `json:"vunit_rate" xml:"VunitRate"`
}

type ValuteOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}
