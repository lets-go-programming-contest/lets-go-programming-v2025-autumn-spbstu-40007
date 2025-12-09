package main

import (
	"fmt"
	"log"

	"github.com/ksuah/task-8/pkg/config"
)

func main() {
	cfg, err := config.Load(config.СonfigFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
