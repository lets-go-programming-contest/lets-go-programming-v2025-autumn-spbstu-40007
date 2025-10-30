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

var ErrUnsupportedCharset = errors.New("неподдерживаемая кодировка")

type CurrencyList struct {
	XMLName xml.Name       `xml:"ValCurs"`
	Items   []CurrencyItem `xml:"Valute"`
}

type CurrencyItem struct {
	XMLName      xml.Name `xml:"Valute"`
	NumberCode   string   `xml:"NumCode"`
	CurrencyCode string   `xml:"CharCode"`
	RawValue     string   `xml:"Value"`
	NumericValue float64  `xml:"-"`
}

type ProcessedCurrency struct {
	NumberCode   int     `json:"num_code"  xml:"NumCode"  yaml:"num_code"`
	CurrencyCode string  `json:"char_code" xml:"CharCode" yaml:"char_code"`
	NumericValue float64 `json:"value"     xml:"Value"    yaml:"value"`
}

type ProcessedCurrencyList struct {
	XMLName xml.Name            `xml:"ValCurs"`
	Items   []ProcessedCurrency `xml:"Valute"`
}

type CurrencySorter []CurrencyItem

func (cs CurrencySorter) Count() int            { return len(cs) }
func (cs CurrencySorter) Exchange(i, j int)     { cs[i], cs[j] = cs[j], cs[i] }
func (cs CurrencySorter) Compare(i, j int) bool { return cs[i].NumericValue > cs[j].NumericValue }

func ProcessXMLFile(filename string) []CurrencyItem {
	fileContent, readErr := readXMLFile(filename)
	if readErr != nil {
		handleFileReadError(filename, readErr)
	}

	currencyData, decodeErr := parseXMLData(fileContent)
	if decodeErr != nil {
		handleXMLParseError(decodeErr)
	}

	return processCurrencyValues(currencyData.Items)
}

func readXMLFile(filename string) ([]byte, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла '%s': %w", filename, err)
	}

	return content, nil
}

func parseXMLData(content []byte) (CurrencyList, error) {
	xmlReader := strings.NewReader(string(content))
	parser := xml.NewDecoder(xmlReader)

	parser.CharsetReader = createCharsetConverter()

	var currencyData CurrencyList
	if err := parser.Decode(&currencyData); err != nil {

		return CurrencyList{}, fmt.Errorf("ошибка декодирования XML: %w", err)
	}

	return currencyData, nil
}

func createCharsetConverter() func(string, io.Reader) (io.Reader, error) {

	return func(encoding string, reader io.Reader) (io.Reader, error) {
		if encoding == "windows-1251" {

			return charmap.Windows1251.NewDecoder().Reader(reader), nil
		}

		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, encoding)
	}
}

func processCurrencyValues(items []CurrencyItem) []CurrencyItem {
	processedItems := make([]CurrencyItem, 0, len(items))

	for _, item := range items {
		processedItem := convertCurrencyValue(item)
		processedItems = append(processedItems, processedItem)
	}

	return processedItems
}

func convertCurrencyValue(item CurrencyItem) CurrencyItem {
	normalizedValue := strings.ReplaceAll(item.RawValue, ",", ".")

	convertedValue, err := strconv.ParseFloat(normalizedValue, 64)
	if err != nil {
		handleValueConversionError(item.RawValue, err)
	}

	item.NumericValue = convertedValue

	return item
}

func (item CurrencyItem) ConvertToOutputFormat() ProcessedCurrency {
	codeValue := convertNumberCode(item.NumberCode)

	return ProcessedCurrency{
		NumberCode:   codeValue,
		CurrencyCode: item.CurrencyCode,
		NumericValue: item.NumericValue,
	}
}

func convertNumberCode(code string) int {
	if code == "" {

		return 0
	}

	convertedCode, err := strconv.Atoi(code)
	if err != nil {
		handleCodeConversionError(code, err)
	}

	return convertedCode
}

func handleFileReadError(filename string, err error) {
	fmt.Printf("Ошибка чтения исходного файла '%s': %v\n", filename, err)
	panic(fmt.Errorf("невозможно прочитать исходный XML файл: %w", err))
}

func handleXMLParseError(err error) {
	fmt.Printf("Ошибка декодирования XML данных: %v\n", err)
	panic(fmt.Errorf("ошибка обработки XML данных: %w", err))
}

func handleValueConversionError(value string, err error) {
	fmt.Printf("Ошибка преобразования значения '%s' в число: %v\n", value, err)
	panic(fmt.Errorf("некорректное значение валюты: %s: %w", value, err))
}

func handleCodeConversionError(code string, err error) {
	fmt.Printf("Ошибка преобразования NumCode '%s' в целое число: %v\n", code, err)
	panic(fmt.Errorf("некорректный NumCode: %s: %w", code, err))
}
