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
	parseAndPrint(configData)
}

func parseAndPrint(data string) {
	lines := strings.Split(data, "\n")
	var env, level string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "environment:") {
			env = strings.TrimSpace(strings.TrimPrefix(line, "environment:"))
		}
		if strings.HasPrefix(line, "log_level:") {
			level = strings.TrimSpace(strings.TrimPrefix(line, "log_level:"))
		}
	}

	fmt.Printf("%s %s\n", env, level)
}
