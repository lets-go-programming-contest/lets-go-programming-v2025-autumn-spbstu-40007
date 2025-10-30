package main

import (
	"flag"
	"fmt"

	"currency-converter/internal/data"
	"currency-converter/internal/output"
)

func main() {
	dir := flag.String("config", "", "путь к директории с входным YAML")
	flag.Parse()

	if *dir == "" {
		panic("Не указан путь к входным данным (используйте --config)")
	}

	currencies, err := data.LoadAndSortCurrencies(*dir)
	if err != nil {
		panic(fmt.Sprintf("Ошибка обработки данных: %v", err))
	}

	output.PrintCurrencies(currencies)
}
