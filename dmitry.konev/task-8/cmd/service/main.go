package main

import (
	"fmt"
	"log"

	"github.com/DichSwitch/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Environment, cfg.LogLevel)
}
