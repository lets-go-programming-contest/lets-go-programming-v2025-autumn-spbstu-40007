package config

import (
	"os"

	yaml "github.com/goccy/go-yaml"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &Config{}
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
