package config

import _ "embed"

var prodConfigBytes []byte

func init() {
	EmbeddedConfig = prodConfigBytes
}