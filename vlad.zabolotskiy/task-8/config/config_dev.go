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
	yaml.Unmarshal(devConfig, &cfg)
	return cfg
}
