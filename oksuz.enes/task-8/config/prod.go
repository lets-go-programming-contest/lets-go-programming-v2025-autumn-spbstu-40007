//go:build !dev

package config

import (
	_ "embed"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var configData []byte

func init() {
	if err := yaml.Unmarshal(configData, &Cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}
