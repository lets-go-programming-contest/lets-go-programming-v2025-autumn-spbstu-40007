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

func New(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var config Config

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return &config, nil
}
