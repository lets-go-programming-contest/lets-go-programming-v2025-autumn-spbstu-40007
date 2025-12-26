package main

import (
	"fmt"
	"log"
)

func main() {
	cfg, err := Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}