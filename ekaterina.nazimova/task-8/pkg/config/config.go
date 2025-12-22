package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func GetConfig() Config {
	var cfg Config

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("Error parsing yaml: %v\n", err)
	}

	return cfg
}
