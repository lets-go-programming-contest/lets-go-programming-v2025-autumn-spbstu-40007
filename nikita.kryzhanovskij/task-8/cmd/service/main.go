package main

import (
	"fmt"
	"strings"

	"task-8/pkg/config"
)

func main() {
	configData := config.GetConfig()

	lines := strings.Split(configData, "\n")
	var environment string

	for _, line := range lines {
		if strings.HasPrefix(line, "environment:") {
			environment = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	fmt.Printf("Loaded config from: pkg/config\n")
	fmt.Printf("Current Environment: %s\n", environment)
}
