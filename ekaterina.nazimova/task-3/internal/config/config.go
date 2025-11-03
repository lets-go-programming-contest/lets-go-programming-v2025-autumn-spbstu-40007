package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const dirPermissions = 0o755

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) (*Config, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	cfg := &Config{
		InputFile:  "",
		OutputFile: "",
	}

	err = yaml.Unmarshal(fileData, cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling config data: %w", err)
	}

	outputDirectory := filepath.Dir(cfg.OutputFile)
	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(outputDirectory, dirPermissions)
		if err != nil {
			return nil, fmt.Errorf("creating output directory %q: %w", outputDirectory, err)
		}
	}

	return cfg, nil
}
