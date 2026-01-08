package main

import (
	"fmt"

	"task-8/config"
)

func main() {
	conf := config.GetConfig()

	fmt.Print(conf.Environment, " ", conf.LogLevel)
}
