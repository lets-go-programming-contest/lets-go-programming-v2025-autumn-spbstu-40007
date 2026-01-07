package main

import (
	"fmt"
	"task-8/config"
)

func main() {
	conf := config.Get()
	fmt.Printf("%s %s\n", conf.Environment, conf.LogLevel)
}
