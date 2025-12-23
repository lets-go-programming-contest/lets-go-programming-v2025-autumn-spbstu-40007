//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var prodFile string

func Load() Config {
	_ = prodFile

	return Config{
		Environment: "prod",
		LogLevel:    "error",
	}
}
