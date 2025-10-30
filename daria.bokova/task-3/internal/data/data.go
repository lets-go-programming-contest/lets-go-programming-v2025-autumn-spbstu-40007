package data

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
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
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	// Создаем декодер с поддержкой windows-1251
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования XML: %v", err)
	}

	var currencies []Currency
	for _, valute := range valCurs.Valutes {
		// Преобразуем значение (заменяем запятую на точку для парсинга)
		valueStr := strings.Replace(valute.Value, ",", ".", -1)

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования значения валюты '%s': %v", valute.Value, err)
		}

		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования числового кода валюты '%s': %v", valute.NumCode, err)
		}

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	// Сортируем по убыванию значения
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies, nil
}
