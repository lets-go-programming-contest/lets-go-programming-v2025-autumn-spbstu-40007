package main

import (
	"flag"

	"github.com/tntkatz/task-3/internal/processor"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./configs/config.yaml", "Write a path to config")
	flag.Parse()

	err := processor.Run(configPath)
	if err != nil {
		panic(err)
	}
}
