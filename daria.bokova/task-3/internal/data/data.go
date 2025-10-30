package data

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type CurrencyData struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	XMLName    xml.Name `xml:"Valute"`
	NumberCode string   `xml:"NumCode"`
	LetterCode string   `xml:"CharCode"`
	RawValue   string   `xml:"Value"`
	ValueFloat float64  `xml:"-"`
}

type ProcessedCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type CurrencyCollection struct {
	XMLName    xml.Name            `xml:"ValCurs"`
	Currencies []ProcessedCurrency `xml:"Valute"`
}

type ByValueDesc []Currency

func (arr ByValueDesc) Len() int           { return len(arr) }
func (arr ByValueDesc) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }
func (arr ByValueDesc) Less(i, j int) bool { return arr[i].ValueFloat > arr[j].ValueFloat }

func ParseCurrencyFile(filename string) []Currency {
	fileContent, readErr := os.ReadFile(filename)
	if readErr != nil {
		log.Printf("Source file reading error '%s': %v", filename, readErr)
		panic("source XML file reading failed: " + readErr.Error())
	}

	xmlReader := xml.NewDecoder(strings.NewReader(string(fileContent)))

	xmlReader.CharsetReader = func(encoding string, reader io.Reader) (io.Reader, error) {
		if encoding == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(reader), nil
		}
		return nil, io.EOF
	}

	var currencyData CurrencyData
	if decodeErr := xmlReader.Decode(&currencyData); decodeErr != nil {
		log.Printf("XML data decoding error: %v", decodeErr)
		panic("XML data decoding failed: " + decodeErr.Error())
	}

	processedData := make([]Currency, 0, len(currencyData.Currencies))

	for _, item := range currencyData.Currencies {
		normalizedValue := strings.ReplaceAll(item.RawValue, ",", ".")

		if parsedValue, convErr := strconv.ParseFloat(normalizedValue, 64); convErr != nil {
			log.Printf("Value conversion error '%s': %v", item.RawValue, convErr)
			panic("invalid currency value: " + item.RawValue + ": " + convErr.Error())
		} else {
			item.ValueFloat = parsedValue
		}

		processedData = append(processedData, item)
	}

	return processedData
}

func (c Currency) ConvertToProcessed() ProcessedCurrency {
	codeNum := 0

	if c.NumberCode != "" {
		if parsedCode, err := strconv.Atoi(c.NumberCode); err != nil {
			log.Printf("Number code conversion error '%s': %v", c.NumberCode, err)
			panic("invalid number code: " + c.NumberCode + ": " + err.Error())
		} else {
			codeNum = parsedCode
		}
	}

	return ProcessedCurrency{
		NumCode:  codeNum,
		CharCode: c.LetterCode,
		Value:    c.ValueFloat,
	}
}
