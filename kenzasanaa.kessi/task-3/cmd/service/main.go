package main

import (
	"flag"

	"kenzasanaa.kessi/task-3/internal"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		panic("--config flag is required")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	err = writer.Run(cfg)
	if err != nil {
		panic(err)
	}
}
