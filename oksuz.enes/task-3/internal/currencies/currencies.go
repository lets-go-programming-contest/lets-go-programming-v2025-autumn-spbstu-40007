package currencies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
)

type Currency struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
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
	if err := xml.Unmarshal(data, &valCurs); err != nil {
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
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create json file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(list); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}
