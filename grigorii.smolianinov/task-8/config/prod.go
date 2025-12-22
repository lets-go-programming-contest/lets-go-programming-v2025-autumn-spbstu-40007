//go:build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var rawConfig []byte

func Load() Config {
	var cfg Config
	_ = yaml.Unmarshal(rawConfig, &cfg)

	return cfg
}
