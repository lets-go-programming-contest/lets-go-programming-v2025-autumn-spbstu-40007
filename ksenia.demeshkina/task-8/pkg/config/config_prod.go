//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var prodYAML []byte

func init() {
	_ = loadYAML(prodYAML)
}
