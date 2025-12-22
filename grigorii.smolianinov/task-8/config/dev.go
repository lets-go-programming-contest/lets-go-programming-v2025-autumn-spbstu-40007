//go:build dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var rawConfig []byte

func init() {
	_ = yaml.Unmarshal(rawConfig, &Current)
}
