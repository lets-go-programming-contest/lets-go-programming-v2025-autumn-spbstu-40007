package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	DevConfig  = []byte(`environment: dev\nlog_level: debug`)
	ProdConfig = []byte(`environment: prod\nlog_level: error`)
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func loadConfig(data []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		
		return nil, err
	}

	return cfg, nil
}

func Load() (*Config, error) {
	cfgData := ProdConfig
	if os.Getenv("GO_ENV") == "dev" {
		cfgData = DevConfig
	}

	return loadConfig(cfgData)
}
