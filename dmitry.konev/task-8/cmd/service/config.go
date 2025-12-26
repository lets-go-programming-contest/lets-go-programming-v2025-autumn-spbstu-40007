package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func devConfig() []byte {
	return []byte(`environment: dev
log_level: debug
`)
}

func prodConfig() []byte {
	return []byte(`environment: prod
log_level: error
`)
}

func loadConfig(data []byte) (*Config, error) {
	cfg := &Config{
		Environment: "",
		LogLevel:    "",
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("yaml unmarshal config: %w", err)
	}

	return cfg, nil
}

func Load() (*Config, error) {
	cfgData := prodConfig()
	if os.Getenv("GO_ENV") == "dev" {
		cfgData = devConfig()
	}

	return loadConfig(cfgData)
}
