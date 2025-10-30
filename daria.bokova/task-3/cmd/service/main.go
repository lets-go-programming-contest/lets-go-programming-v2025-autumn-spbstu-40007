package main

import (
	"flag"
	"fmt"

	"currency-converter/internal/config"
	"currency-converter/internal/data"
	"currency-converter/internal/output"
)

func main() {
	configPath := flag.String("config", "", "Путь до конфигурационного файла")
	flag.Parse()

	if *configPath == "" {
		panic("Не указан путь до конфигурационного файла. Используйте флаг --config")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Ошибка загрузки конфигурации: %v", err))
	}

	currencies, err := data.LoadAndSortCurrencies(cfg.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Ошибка обработки данных: %v", err))
	}

	err = output.SaveCurrencies(currencies, cfg.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("Ошибка сохранения результатов: %v", err))
	}

	fmt.Printf("Данные успешно обработаны и сохранены в %s\n", cfg.OutputFile)
}
