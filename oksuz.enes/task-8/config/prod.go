//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var prodConfigData string

func init() {
	Environment = "prod"
	LogLevel = "error"
}
