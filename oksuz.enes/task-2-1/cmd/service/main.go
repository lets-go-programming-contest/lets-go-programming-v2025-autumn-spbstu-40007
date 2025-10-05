package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var room, workers int
	if _, err := fmt.Fscan(in, &room, &workers); err != nil {

		return
	}

	minTemp := 15
	maxTemp := 30

	for i := 0; i < workers; i++ {
		var op string
		var t int

		if _, err := fmt.Fscan(in, &op, &t); err != nil {

			return
		}

		if op == ">=" {
			if t > minTemp {
				minTemp = t
			}
		} else if op == "<=" {
			if t < maxTemp {
				maxTemp = t
			}
		}
	}

	if minTemp > maxTemp {
		fmt.Println(-1)

		return
	}

	fmt.Println(minTemp)
}
