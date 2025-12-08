package config

import (
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func New(path string) (*Config, error) {
	configContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	config := &Config{
		InputFile:  "",
		OutputFile: "",
	}
	if err = yaml.Unmarshal(configContent, config); err != nil {
		return nil, fmt.Errorf("config: did not find expected key %w", err)
	}

	return config, nil
}
