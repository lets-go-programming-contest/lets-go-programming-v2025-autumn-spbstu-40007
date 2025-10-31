package structs

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type FloatVal float64

type Valute struct {
	NumCode  int      `json:"num_code"  xml:"NumCode"`
	CharCode string   `json:"char_code" xml:"CharCode"`
	Value    FloatVal `json:"value"     xml:"Value"`
}

type Valutes struct {
	Valutes []Valute `xml:"Valute"`
}

func (floatVal *FloatVal) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := decoder.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("error decoding xml: %w", err)
	}

	s = strings.ReplaceAll(s, ",", ".")

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("error converting to float: %w", err)
	}

	*floatVal = FloatVal(val)
	return nil
}
