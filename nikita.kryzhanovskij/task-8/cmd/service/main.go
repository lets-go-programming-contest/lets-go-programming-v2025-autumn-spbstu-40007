package main

import (
	"fmt"
	"strings"

	"task-8/pkg/config"
)

func main() {
	configData := config.GetConfig()

	lines := strings.Split(configData, "\n")
	var env, level string

	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
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

	fmt.Printf("%s %s\n", env, level)
}
