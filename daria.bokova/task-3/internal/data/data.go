package data

import (
	"fmt"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

type Currency struct {
	NumCode  int     `yaml:"num"`
	CharCode string  `yaml:"code"`
	Value    float64 `yaml:"value"`
}

type Input struct {
	Currencies []Currency `yaml:"currencies"`
}

func LoadAndSortCurrencies(path string) ([]Currency, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	var input Input
	err = yaml.Unmarshal(bytes, &input)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %v", err)
	}

	if len(input.Currencies) == 0 {
		return nil, fmt.Errorf("файл не содержит валют")
	}

	// сортировка: по убыванию Value, при равенстве — по возрастанию NumCode
	slices.SortFunc(input.Currencies, func(a, b Currency) int {
		if a.Value > b.Value {
			return -1
		}
		if a.Value < b.Value {
			return 1
		}
		if a.NumCode < b.NumCode {
			return -1
		}
		if a.NumCode > b.NumCode {
			return 1
		}
		return 0
	})

	return input.Currencies, nil
}
