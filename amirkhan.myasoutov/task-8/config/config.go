package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type AppOptions struct {
	Stage string `yaml:"environment"`
	Level string `yaml:"log_level"`
}

func (o AppOptions) String() string {
	return fmt.Sprintf("Environment: %s, LogLevel: %s", o.Stage, o.Level)
}

func decodeYaml(blob []byte) (*AppOptions, error) {
	var opt AppOptions

	if err := yaml.Unmarshal(blob, &opt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return &opt, nil
}