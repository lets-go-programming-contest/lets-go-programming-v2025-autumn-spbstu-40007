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

func LoadConfig(path string) Config {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		panic(fmt.Errorf("failed to read configuration file '%s': %w", path, readErr))
	}

	var cfg Config

	decodeErr := yaml.Unmarshal(data, &cfg)
	if decodeErr != nil {
		panic(fmt.Errorf("failed to decode YAML configuration: %w", decodeErr))
	}

	return cfg
}
