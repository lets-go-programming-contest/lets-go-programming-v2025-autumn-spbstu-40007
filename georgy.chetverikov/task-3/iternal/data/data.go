package data

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Valutes []Valute

type Valute struct {
	NumCode  int         `xml:"NumCode" json:"num_code"`
	CharCode string      `xml:"CharCode" json:"char_code"`
	Nominal  int         `xml:"Nominal" json:"nominal"`
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

func ParseXML(data []byte) (Valutes, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}

	type tempValCourse struct {
		Valutes []Valute `xml:"Valute"`
	}

	var temp tempValCourse
	err := decoder.Decode(&temp)
	if err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	return temp.Valutes, nil
}
