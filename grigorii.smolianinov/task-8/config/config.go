package config

import "fmt"

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func PrintInfo(cfg Config) {
	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
