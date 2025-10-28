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

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to read input file %q: %w", path, err)
	}

	config := new(Config)

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("Unable to unmarshall data: %w", err)
	}

	return config, nil
}
