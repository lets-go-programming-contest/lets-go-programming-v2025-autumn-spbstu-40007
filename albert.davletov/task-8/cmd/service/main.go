package main

import (
	"fmt"

	"github.com/treadwave/task-8/config"
)

func main() {
	conf := config.GetConfig()

	fmt.Print(conf.Environment, " ", conf.LogLevel)
}
