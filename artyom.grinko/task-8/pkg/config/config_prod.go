//go:build prod

package config

import _ "embed"

//go:embed prod.yaml
var configContents string
