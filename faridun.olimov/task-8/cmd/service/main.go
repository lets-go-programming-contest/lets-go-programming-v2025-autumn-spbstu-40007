//nolint:gofumpt,gofmt
package main

import (
	"fmt"
	"log"

	"task-8/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфиг %v", err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
