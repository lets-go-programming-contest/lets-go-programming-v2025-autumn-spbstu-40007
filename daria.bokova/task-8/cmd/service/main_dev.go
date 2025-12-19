//go:build dev

package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed dev.yaml
var configData string

func main() {
	lines := strings.Split(configData, "\n")
	env, level := "", ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "environment:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				env = strings.TrimSpace(parts[1])
				env = strings.Trim(env, "\"")
			}
		}
		if strings.HasPrefix(line, "log_level:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				level = strings.TrimSpace(parts[1])
				level = strings.Trim(level, "\"")
			}
		}
	}

	fmt.Printf("%s %s", env, level)
}
