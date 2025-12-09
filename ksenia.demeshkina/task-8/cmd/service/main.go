package main

import (
	"fmt"

	"github.com/ksuah/task-8/pkg/config"
)

func main() {
	c := config.Get()
	fmt.Printf("%s %s\n", c.Environment, c.LogLevel)
}
