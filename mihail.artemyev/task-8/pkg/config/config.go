package config

import "gopkg.in/yaml.v3"

type Configuration struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (Configuration, error) {
	var configuration Configuration

	if err := yaml.Unmarshal(configBytes, &configuration); err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}
