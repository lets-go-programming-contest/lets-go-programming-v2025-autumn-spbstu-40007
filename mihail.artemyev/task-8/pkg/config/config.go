package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func ParseConfiguration(rawData []byte) (*Configuration, error) {
	configuration := &Configuration{}

	parseError := yaml.Unmarshal(rawData, configuration)
	if parseError != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", parseError)
	}

	return configuration, nil
}
