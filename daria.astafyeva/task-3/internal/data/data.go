package data

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string  `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	ValueStr string  `xml:"Value"`
	Value    float64 `xml:"-"`
}

type ResultValute struct {
	NumCode  int     `json:"num_code" xml:"NumCode" yaml:"num_code"`
	CharCode string  `json:"char_code" xml:"CharCode" yaml:"char_code"`
	Value    float64 `json:"value" xml:"Value" yaml:"value"`
}

type ResultValutes struct {
	XMLName xml.Name       `xml:"ValCurs"`
	Valutes []ResultValute `xml:"Valute"`
}

type ByValueDesc []Valute

func (a ByValueDesc) Len() int           { return len(a) }
func (a ByValueDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValueDesc) Less(i, j int) bool { return a[i].Value > a[j].Value }

func DecodeXMLData(filePath string) []Valute {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to read source XML file '%s': %w", filePath, err))
	}

	reader := bytes.NewReader(fileData)
	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}

	var valCurs ValCurs
	if decodeErr := decoder.Decode(&valCurs); decodeErr != nil {
		panic(fmt.Errorf("XML data decoding error: %w", decodeErr))
	}

	processedValutes := make([]Valute, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		valueStr := strings.Replace(valute.ValueStr, ",", ".", 1)

		value, parseErr := strconv.ParseFloat(valueStr, 64)
		if parseErr != nil {
			panic(fmt.Errorf("invalid currency value '%s': %w", valute.ValueStr, parseErr))
		}

		valute.Value = value
		processedValutes = append(processedValutes, valute)
	}
	return processedValutes
}

func (v Valute) ToResultValute() ResultValute {
	numCode := 0

	if v.NumCode != "" {
		nc, err := strconv.Atoi(v.NumCode)
		if err != nil {
			panic(fmt.Errorf("invalid NumCode '%s': %w", v.NumCode, err))
		}
		numCode = nc
	}

	return ResultValute{
		NumCode:  numCode,
		CharCode: v.CharCode,
		Value:    v.Value,
	}
}
