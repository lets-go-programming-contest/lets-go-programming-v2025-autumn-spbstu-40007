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
	var valString string
	if err := decoder.DecodeElement(&valString, &start); err != nil {
		return fmt.Errorf("error decoding xml: %w", err)
	}

	valString = strings.ReplaceAll(valString, ",", ".")

	val, err := strconv.ParseFloat(valString, 64)
	if err != nil {
		return fmt.Errorf("error converting to float: %w", err)
	}

	*floatVal = FloatVal(val)

	return nil
}
