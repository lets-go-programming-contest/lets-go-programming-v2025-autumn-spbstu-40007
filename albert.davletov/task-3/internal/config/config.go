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

func YamlDecoder(filepath string) (Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file: %w", err)
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	var cfg Config
	err = decoder.Decode(&cfg)

	if err != nil {
		return Config{}, fmt.Errorf("error decoding YAML: %w", err)
	}

	return cfg, nil
}
