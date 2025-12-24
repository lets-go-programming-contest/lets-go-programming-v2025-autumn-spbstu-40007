//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var rawBytes []byte

func Fetch() (*AppOptions, error) {
	return decodeYaml(rawBytes)
}