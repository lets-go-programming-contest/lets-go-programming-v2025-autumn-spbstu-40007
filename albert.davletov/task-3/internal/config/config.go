package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func YamlDecoder(filepath string) (Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file: %w", err)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Printf("error closing file: %v\n", closeErr)
		}
	}()

	decoder := yaml.NewDecoder(file)

	var cfg Config

	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error decoding YAML: %w", err)
	}

	return cfg, nil
}
