package main

import (
	"fmt"

	conf "task-8/config"
)

func main() {
	c := conf.Current()

	fmt.Println(c.Env, c.Log)
}
