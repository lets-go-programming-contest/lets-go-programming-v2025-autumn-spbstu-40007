package xmlparser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	ValueStr string  `xml:"Value"`
	Value    float64 `xml:"-"`
}

func LoadXML(path string) []Valute {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("Cannot open XML file: %v", err)
	}

	reader := bytes.NewReader(data)

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return input, nil
	}

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		log.Panicf("Cannot parse XML: %v", err)
	}

	for i := range valCurs.Valutes {
		valCurs.Valutes[i].Value = parseValue(valCurs.Valutes[i].ValueStr)
	}

	return valCurs.Valutes
}

func parseValue(s string) float64 {
	s = strings.Replace(s, ",", ".", 1)

	var val float64
	
	_, err := fmt.Sscanf(s, "%f", &val)

	if err != nil {
		log.Panicf("Invalid number format in XML: %s", s)
	}

	return val
}
