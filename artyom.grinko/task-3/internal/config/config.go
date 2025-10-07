package config

import (
	"errors"
	"os"

	"task-3/internal/die"

	yaml "github.com/goccy/go-yaml"
)

var errDidNotFindExpectedKey = errors.New("config: did not find expected key")

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			die.Die(err)
		}
	}()

	result := &Config{} //nolint:exhaustruct
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(result); err != nil {
		return nil, errDidNotFindExpectedKey
	}

	return result, nil
}
