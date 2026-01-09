//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var configData string

func GetConfig() string {
	return configData
}
