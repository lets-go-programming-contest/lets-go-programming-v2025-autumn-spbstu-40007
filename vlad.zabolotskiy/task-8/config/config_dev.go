//go:build dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfig []byte

func getConfig() Config {
	var cfg Config

	if err := yaml.Unmarshal(devConfig, &cfg); err != nil {
		return Config{
			Environment: "dev",
			LogLevel:    "debug",
		}
	}

	return cfg
}
