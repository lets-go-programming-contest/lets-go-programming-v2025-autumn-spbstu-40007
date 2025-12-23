package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func loadConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {

		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func Load() (*Config, error) {
	cfgData := ProdConfig 
	if os.Getenv("GO_ENV") == "dev" {
		cfgData = DevConfig
	}
	
	return loadConfig(cfgData)
}
