//go:build !dev

package settings

import _ "embed"

//go:embed prod.yaml
var prodConfigData []byte

func init() {
	GetConfig = func() []byte {
		return prodConfigData
	}
}
