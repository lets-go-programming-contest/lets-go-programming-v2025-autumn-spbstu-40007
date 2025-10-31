package data

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type xmlCurrency struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float32 `json:"value"`
}

type ValCurs struct {
	XMLName    xml.Name      `xml:"ValCurs"`
	Currencies []xmlCurrency `xml:"Valute"`
}

func LoadFromXML(path string) ([]Currency, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(string(data))
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)

	if err != nil {
		return nil, err
	}

	var currencies []Currency
	for _, xmlCurr := range valCurs.Currencies {
		valueStr := strings.Replace(xmlCurr.Value, ",", ".", -1)

		value, err := strconv.ParseFloat(valueStr, 32)

		if err != nil {
			return nil, err
		}

		currencies = append(currencies, Currency{
			NumCode:  xmlCurr.NumCode,
			CharCode: xmlCurr.CharCode,
			Value:    float32(value),
		})
	}

	return currencies, nil
}

func SortByValue(currencies []Currency) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
}

func SaveToJSON(currencies []Currency, path string) error {
	folderPath := filepath.Dir(path)

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	jsonData, err := json.MarshalIndent(currencies, "", "  ")

	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	return err
}
