//go:build dev

package config

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed dev.yaml
var devConfig string

type Config struct {
	Environment string
	LogLevel    string
}

func parseYAML(content string) (*Config, error) {
	cfg := &Config{}
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

		switch key {
		case "environment":
			cfg.Environment = value
		case "log_level":
			cfg.LogLevel = value
		}
	}

	return cfg, nil
}

func Load() (*Config, error) {
	return parseYAML(devConfig)
}

func Print() {
	cfg, err := Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
