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
	err := yaml.Unmarshal(configData, &cfg)
	if err != nil {
		// Обработка ошибки - просто выходим
		return
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
