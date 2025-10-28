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

func LoadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("Failed to read config file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	
	if err != nil {
		log.Panicf("Failed to parse YAML: %v", err)
	}

	return cfg
}
