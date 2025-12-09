package config

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment" json:"environment"`
	LogLevel    string `yaml:"log_level" json:"log_level"`
}

var cfg Config

func Get() Config {
	return cfg
}

func loadYAML(data []byte) error {
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		if err2 := json.Unmarshal(data, &cfg); err2 != nil {
			return err
		}
	}
	return nil
}
