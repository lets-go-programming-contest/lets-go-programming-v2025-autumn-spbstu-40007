//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var configData string

func GetConfig() string {
	return configData
}
