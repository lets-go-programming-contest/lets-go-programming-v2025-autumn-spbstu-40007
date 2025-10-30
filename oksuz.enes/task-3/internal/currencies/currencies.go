package currencies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/html/charset"
)

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type ValCurs struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}

type CurrencyService struct{}

func NewCurrencyService() *CurrencyService {
	return &CurrencyService{}
}

func (s *CurrencyService) ParseXML(data []byte) ([]Currency, error) {
	var valCurs ValCurs
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to parse xml: %w", err)
	}

	return valCurs.Currencies, nil
}

func (s *CurrencyService) SortByValue(list []Currency) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Value > list[j].Value
	})
}

func (s *CurrencyService) SaveToJSON(path string, list []Currency) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create json file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Println("close error:", cerr)
		}
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(list); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	return nil
}
