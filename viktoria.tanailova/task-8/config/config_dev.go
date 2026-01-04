//go:build dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devData []byte

func load() Config {
	var cfg Config
	_ = yaml.Unmarshal(devData, &cfg)
	return cfg
}
