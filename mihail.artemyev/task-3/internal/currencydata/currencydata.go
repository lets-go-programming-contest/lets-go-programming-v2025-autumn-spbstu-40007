package currencydata

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (currencyValue *CurrencyValue) UnmarshalXML(
	xmlDecoder *xml.Decoder,
	startElement xml.StartElement,
) error {
	var xmlStringContent string

	decodingError := xmlDecoder.DecodeElement(&xmlStringContent, &startElement)
	if decodingError != nil {
		return fmt.Errorf("decoding XML element: %w", decodingError)
	}

	normalizedStringValue := strings.Replace(xmlStringContent, ",", ".", 1)

	parsedFloatValue, parsingError := strconv.ParseFloat(normalizedStringValue, 64)
	if parsingError != nil {
		return fmt.Errorf("parsing float value from string %q: %w", xmlStringContent, parsingError)
	}

	*currencyValue = CurrencyValue(parsedFloatValue)

	return nil
}

type CurrencyExchange struct {
	NumericCode int              `json:"num_code"  xml:"NumCode"`
	CharCode    string           `json:"char_code" xml:"CharCode"`
	ExchangeRate CurrencyValue   `json:"value"     xml:"Value"`
}

type ExchangeRateList struct {
	ExchangeDate string                `xml:"Date,attr"`
	MarketName   string                `xml:"name,attr"`
	CurrencyList []CurrencyExchange    `xml:"Valute"`
}

type CurrencyCollection []CurrencyExchange

func (collection CurrencyCollection) Len() int {
	return len(collection)
}

func (collection CurrencyCollection) Swap(firstIndex, secondIndex int) {
	collection[firstIndex], collection[secondIndex] = collection[secondIndex], collection[firstIndex]
}

func (collection CurrencyCollection) Less(firstIndex, secondIndex int) bool {
	return collection[firstIndex].ExchangeRate > collection[secondIndex].ExchangeRate
}
