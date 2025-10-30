package config

import (
	"log"
	"os"

	yaml "github.com/goccy/go-yaml"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("Failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Panicf("Failed to parse YAML: %v", err)
	}

	return &cfg
}
