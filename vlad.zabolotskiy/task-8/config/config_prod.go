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
	yaml.Unmarshal(prodConfig, &cfg)
	return cfg
}
