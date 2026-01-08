package config

import _ "embed"

var devConfigBytes []byte

func init() {
	EmbeddedConfig = devConfigBytes
}