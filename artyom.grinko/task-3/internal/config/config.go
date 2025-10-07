package config

import (
	"errors"
	"os"

	yaml "github.com/goccy/go-yaml"
)

var didNotFindExpectedKey = errors.New("config: did not find expected key")

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
		return nil, didNotFindExpectedKey
	}

	return result, nil
}
