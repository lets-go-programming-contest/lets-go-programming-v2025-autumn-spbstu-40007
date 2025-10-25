package config

import (
	"log"
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
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	outputDir := filepath.Dir(cfg.OutputFile)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Printf("Creating output directory: %s", outputDir)
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
