package data

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCourse struct {
	Date    string   `xml:"Date,attr" json:"date"`
	Valutes []Valute `xml:"Valute" json:"valutes"`
}

type Valute struct {
	NumCode  string      `xml:"NumCode" json:"num_code"`
	CharCode string      `xml:"CharCode" json:"char_code"`
	Nominal  string      `xml:"Nominal" json:"nominal"`
	Value    customFloat `xml:"Value" json:"value"`
}

type customFloat float64

func (f *customFloat) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var XMLdata string

	if err := decoder.DecodeElement(&XMLdata, &start); err != nil {
		return fmt.Errorf("decoding xml-file: %w", err)
	}

	decodedData := strings.Replace(XMLdata, ",", ".", 1)

	parsedData, err := strconv.ParseFloat(decodedData, 64)
	if err != nil {
		return fmt.Errorf("parsing data: %w", err)
	}

	*f = customFloat(parsedData)

	return nil
}
