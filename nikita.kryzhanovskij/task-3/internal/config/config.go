package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"nikita.kryzhanovskij/task-3/internal/models"
)

var ErrInvalidConfig = errors.New("invalid config: input-file and output-file are required")

func Load(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("no such file or directory: %w", err)
		}

		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg models.Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return nil, ErrInvalidConfig
	}

	return &cfg, nil
}
