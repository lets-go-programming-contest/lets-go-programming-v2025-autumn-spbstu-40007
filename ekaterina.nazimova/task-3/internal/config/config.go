package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling config data: %w", err)
	}

	outputDir := filepath.Dir(cfg.OutputFile)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {

		err := os.MkdirAll(outputDir, 0o755)
		if err != nil {
			return nil, fmt.Errorf("creating output directory %q: %w", outputDir, err)
		}
	}

	return cfg, nil
}
