package main

import (
	"fmt"

	conf "task-8/config"
)

func main() {
	c := conf.Current()

	fmt.Print(c.Env, " ", c.Log)
}
