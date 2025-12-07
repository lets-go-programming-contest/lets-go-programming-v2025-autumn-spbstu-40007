package config

import (
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() Config {
	var cfg Config

	err := yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка чтения конфига: %v", err)
	}

	return cfg
}
