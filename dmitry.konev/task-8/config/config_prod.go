//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var configData []byte

func Load() (*Config, error) {
	return loadConfig(configData)
}