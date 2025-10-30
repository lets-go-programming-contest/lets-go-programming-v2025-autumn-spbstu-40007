package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурационного файла: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования YAML: %v", err)
	}

	if config.InputFile == "" {
		return nil, fmt.Errorf("поле input-file не может быть пустым")
	}

	if config.OutputFile == "" {
		return nil, fmt.Errorf("поле output-file не может быть пустым")
	}

	return &config, nil
}
