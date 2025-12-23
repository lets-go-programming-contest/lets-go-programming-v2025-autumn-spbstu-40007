//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var devConfigData string

func init() {
	Environment = "dev"
	LogLevel = "debug"
}
