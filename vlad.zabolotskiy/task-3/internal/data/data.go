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
		return nil, err //nolint:wrapcheck
	}

	reader := strings.NewReader(string(data))
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	var currencies []Currency //nolint:prealloc

	for _, xmlCurr := range valCurs.Currencies {
		valueStr := strings.ReplaceAll(xmlCurr.Value, ",", ".")

		value, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			return nil, err //nolint:wrapcheck
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
	if err := os.MkdirAll(folderPath, 0o755); err != nil { //nolint:mnd
		return err //nolint:wrapcheck
	}

	file, err := os.Create(path)
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer file.Close() //nolint:errcheck

	jsonData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return err //nolint:wrapcheck
	}

	_, err = file.Write(jsonData)

	return err //nolint:nlreturn,wrapcheck
}
