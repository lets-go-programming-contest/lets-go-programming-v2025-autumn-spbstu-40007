package main

import (
	"fmt"
	"log"
	"github.com/itsdasha/task-8/pkg/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile) 
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	fmt.Printf("Environment: %s, LogLevel: %s\n", cfg.Environment, cfg.LogLevel)
}