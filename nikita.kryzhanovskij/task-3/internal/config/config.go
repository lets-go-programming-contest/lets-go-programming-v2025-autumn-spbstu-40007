package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"nikita.kryzhanovskij/task-3/internal/models"
)

func Load(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return nil, fmt.Errorf("invalid config: input-file and output-file are required")
	}

	return &cfg, nil
}
