package config

import (
	_ "embed"
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfig []byte

//go:embed prod.yaml
var prodConfig []byte

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfgData []byte

	cfgData = prodConfig
	if os.Getenv("GO_ENV") == "dev" {
		cfgData = devConfig
	}

	return loadConfig(cfgData)
}

func loadConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {

		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &cfg, nil
}

