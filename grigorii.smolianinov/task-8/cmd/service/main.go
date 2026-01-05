package main

import "grigorii.smolianinov/task-8/config"

func main() {
	cfg := config.Load()

	config.PrintInfo(cfg)
}
