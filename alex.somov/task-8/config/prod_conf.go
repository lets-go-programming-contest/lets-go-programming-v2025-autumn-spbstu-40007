//go:build !dev

package conf

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodRaw []byte

func load() Config {
	var cfg Config
	_ = yaml.Unmarshal(prodRaw, &cfg)
	return cfg
}
