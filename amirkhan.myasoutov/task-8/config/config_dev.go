//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var rawBytes []byte

func Fetch() (*AppOptions, error) {
	return decodeYaml(rawBytes)
}
