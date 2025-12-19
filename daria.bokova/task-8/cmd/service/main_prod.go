//go:build !dev

package main

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var configData []byte

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func main() {
	var cfg Config
	yaml.Unmarshal(configData, &cfg)
	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
