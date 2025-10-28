// работа с конфигурацией YAML
package config

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) Config { // возвращает значение типа Config, с помощью нее получаем параметры из файла конфигурации, вызываем в vain
	data, err := os.ReadFile(path) // читаем весь файл кинфигурации в срез байт data
	if err != nil {
		log.Panicf("Failed to read config file: %v", err)
	}

	var cfg Config                   // переменная cfg типа Config, в нее парсим YAML
	err = yaml.Unmarshal(data, &cfg) // декодируем YAML-байты в структуру go, Unmarshal заполняет структуру по ссылке
	if err != nil {
		log.Panicf("Failed to parse YAML: %v", err)
	}

	return cfg // заполненая структура
}
