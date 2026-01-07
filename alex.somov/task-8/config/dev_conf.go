//go:build dev

package conf

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devRaw []byte

func load() Config {
	var cfg Config
	_ = yaml.Unmarshal(devRaw, &cfg)
	return cfg
}
