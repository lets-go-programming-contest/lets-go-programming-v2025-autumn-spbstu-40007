package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

var ErrKeyNotFound = errors.New("did not find expected key")

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
		return nil, fmt.Errorf("%w: %w", ErrKeyNotFound, err)
	}

	return &cfg, nil
}
