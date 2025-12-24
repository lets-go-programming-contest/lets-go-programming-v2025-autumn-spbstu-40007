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
	return fmt.Sprintf("%s %s", o.Stage, o.Level)
}

func decodeYaml(blob []byte) (*AppOptions, error) {
	var opt AppOptions
	if err := yaml.Unmarshal(blob, &opt); err != nil {
		return nil, fmt.Errorf("decode settings: %w", err)
	}
	return &opt, nil
}
