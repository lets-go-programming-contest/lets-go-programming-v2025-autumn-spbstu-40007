package config

import _ "embed"

//go:embed dev.yaml
var DevConfig []byte

//go:embed prod.yaml
var ProdConfig []byte
