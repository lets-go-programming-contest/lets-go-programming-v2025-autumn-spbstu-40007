package main

import (
	"fmt"
	"strings"

	"task-8/pkg/config"
)

const splitParts = 2

func main() {
	configData := config.GetConfig()
	lines := strings.Split(configData, "\n")

	var env, level string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", splitParts)
		if len(parts) != splitParts {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if key == "environment" {
			env = val
		} else if key == "log_level" {
			level = val
		}
	}

	if env != "" && level != "" {
		fmt.Printf("%s %s", env, level)
	}
}
