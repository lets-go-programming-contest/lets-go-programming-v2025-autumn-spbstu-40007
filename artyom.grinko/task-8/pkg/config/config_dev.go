//go:build !prod

package config

import _ "embed"

//go:embed dev.yaml
var configContents string
