//go:build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodData []byte

func load() Config {
	var cfg Config
	_ = yaml.Unmarshal(prodData, &cfg)
	return cfg
}
