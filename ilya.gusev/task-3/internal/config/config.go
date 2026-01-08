package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Configuration struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(filePath string) (*Configuration, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var cfg Configuration

	if err := yaml.Unmarshal(fileData, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse yaml config: %w", err)
	}

	return &cfg, nil
}
