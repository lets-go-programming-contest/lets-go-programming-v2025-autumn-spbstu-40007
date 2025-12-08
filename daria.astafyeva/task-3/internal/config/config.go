package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	SourceFile string `yaml:"input-file"`
	TargetFile string `yaml:"output-file"`
}

func LoadSettings(path string) Settings {
	rawData, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("cannot read config file '%s': %w", path, err))
	}

	var cfg Settings
	if err := yaml.Unmarshal(rawData, &cfg); err != nil {
		panic(fmt.Errorf("cannot parse YAML config: %w", err))
	}

	return cfg
}
