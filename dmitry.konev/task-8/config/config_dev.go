//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var configData []byte

func Load() (*Config, error) {
	return loadConfig(configData)
}
