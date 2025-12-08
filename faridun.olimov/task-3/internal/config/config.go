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
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading configuration file '%s': %v\n", path, err)
		panic(fmt.Errorf("failed to read configuration file: %w", err))
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("Error decoding YAML configuration: %v\n", err)
		panic(fmt.Errorf("failed to decode YAML configuration: %w", err))
	}

	return cfg
}
