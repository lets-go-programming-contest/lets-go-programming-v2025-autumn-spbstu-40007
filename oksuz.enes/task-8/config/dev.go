//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var devFile string

func Load() Config {
	_ = devFile

	return Config{
		Environment: "dev",
		LogLevel:    "debug",
	}
}
