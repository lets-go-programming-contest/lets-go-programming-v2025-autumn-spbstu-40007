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
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	decoder := yaml.NewDecoder(file)

	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("decode yaml: %w", err)
	}

	if err = file.Close(); err != nil {
		return nil, fmt.Errorf("close file: %w", err)
	}

	return &config, nil
}
