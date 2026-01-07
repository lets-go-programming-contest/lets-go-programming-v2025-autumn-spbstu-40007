//go:build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfig []byte

func getConfig() Config {
	var cfg Config

	if err := yaml.Unmarshal(prodConfig, &cfg); err != nil {
		return Config{
			Environment: "prod",
			LogLevel:    "error",
		}
	}

	return cfg
}
