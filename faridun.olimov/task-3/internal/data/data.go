package data

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

var ErrUnknownCharset = errors.New("unknown charset")

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	ValueStr string   `xml:"Value"`
	Value    float64  `xml:"-"`
}

type ResultValute struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"  yaml:"num_code"`
	CharCode string  `json:"char_code" xml:"CharCode" yaml:"char_code"`
	Value    float64 `json:"value"     xml:"Value"    yaml:"value"`
}

type ResultValutes struct {
	XMLName xml.Name       `xml:"ValCurs"`
	Valutes []ResultValute `xml:"Valute"`
}

type CustomSorter []Valute

func (a CustomSorter) Len() int           { return len(a) }
func (a CustomSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CustomSorter) Less(i, j int) bool { return a[i].Value > a[j].Value }

func DecodeXMLData(filePath string) []Valute {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading source file '%s': %v\n", filePath, err)
		panic(fmt.Errorf("failed to read source XML file: %w", err))
	}

	decoder := xml.NewDecoder(strings.NewReader(string(data)))

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, fmt.Errorf("%w: %s", ErrUnknownCharset, charset)
	}

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		fmt.Printf("Error decoding XML data: %v\n", err)
		panic(fmt.Errorf("XML data decoding error: %w", err))
	}

	processedValutes := make([]Valute, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		valueStr := strings.ReplaceAll(valute.ValueStr, ",", ".")
		value, err := strconv.ParseFloat(valueStr, 64)

		if err != nil {
			fmt.Printf("Error converting value '%s' to float64: %v\n", valute.ValueStr, err)
			panic(fmt.Errorf("invalid currency value: %s: %w", valute.ValueStr, err))
		}

		valute.Value = value
		processedValutes = append(processedValutes, valute)
	}

	return processedValutes
}

func (v Valute) ToResultValute() ResultValute {
	numCode := 0 
	
	if v.NumCode != "" {
		var err error
		numCode, err = strconv.Atoi(v.NumCode)

		if err != nil {
			fmt.Printf("Error converting NumCode '%s' to integer: %v\n", v.NumCode, err)
			panic(fmt.Errorf("invalid NumCode: %s: %w", v.NumCode, err))
		}
	} 

	return ResultValute{
		NumCode:  numCode,
		CharCode: v.CharCode,
		Value:    v.Value,
	}
}
