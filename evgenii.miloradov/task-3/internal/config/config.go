package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(configPath string) (*AppConfig, error) {
	configData, readErr := os.ReadFile(configPath)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read config file: %w", readErr)
	}

	cfg := &AppConfig{}

	parseErr := yaml.Unmarshal(configData, cfg)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse config data: %w", parseErr)
	}

	return cfg, nil
}
