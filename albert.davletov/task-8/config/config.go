package config

import (
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func GetConfig() *Config {
	return loadConfig()
}

func ParseYAML(data []byte) *Config {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("error parsing YAML: %v", err)
	}
	return &cfg
}
