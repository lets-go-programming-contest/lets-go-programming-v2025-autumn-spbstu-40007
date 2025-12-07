//go:build prod

package config

import _ "embed"

//go:embed prod.yaml
var configFile []byte
