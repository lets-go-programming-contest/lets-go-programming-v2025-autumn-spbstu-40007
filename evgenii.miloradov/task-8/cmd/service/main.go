package main

import (
	"fmt"

	config "task-8/config"
)

func main() {
	c := config.Current()

	fmt.Print(c.Env, " ", c.Log)
}
