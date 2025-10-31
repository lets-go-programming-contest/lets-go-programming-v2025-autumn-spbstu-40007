package config

import (
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
		return nil, err //nolint:wrapcheck
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &config, nil
}
