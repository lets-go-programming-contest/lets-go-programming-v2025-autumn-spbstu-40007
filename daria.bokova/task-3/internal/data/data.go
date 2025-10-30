package data

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Value    string   `xml:"Value"`
}

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func LoadAndSortCurrencies(inputFile string) ([]Currency, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	var valCurs ValCurs
	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования XML: %v", err)
	}

	var currencies []Currency
	for _, valute := range valCurs.Valutes {
		// Преобразуем значение (заменяем запятую на точку для парсинга)
		valueStr := valute.Value
		valueStr = replaceCommaWithDot(valueStr)

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования значения валюты: %v", err)
		}

		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования числового кода валюты: %v", err)
		}

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies, nil
}

func replaceCommaWithDot(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r == ',' {
			result[i] = '.'
		} else {
			result[i] = r
		}
	}
	return string(result)
}
