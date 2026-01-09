package main

import (
	"fmt"

	"github.com/se1lzor/task-8/config"
)

func main() {
	conf := config.Get()
	fmt.Printf("%s %s", conf.Environment, conf.LogLevel)
}
