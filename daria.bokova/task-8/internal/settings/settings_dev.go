//go:build dev

package settings

import _ "embed"

//go:embed dev.yaml
var devConfigData []byte

func init() {
	GetConfig = func() []byte {
		return devConfigData
	}
}
