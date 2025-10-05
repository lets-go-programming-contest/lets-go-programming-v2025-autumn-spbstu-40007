package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var room, workers int
	if _, err := fmt.Fscan(reader, &room, &workers); err != nil {
		return
	}

	minTemp := 15
	maxTemp := 30

	for range workers {
		var operator string
		var temp int

		if _, err := fmt.Fscan(reader, &operator, &temp); err != nil {
			return
		}

		if operator == ">=" {
			if temp > minTemp {
				minTemp = temp
			}
		} else if operator == "<=" {
			if temp < maxTemp {
				maxTemp = temp
			}
		}
	}

	if minTemp > maxTemp {
		fmt.Println(-1)
	} else {
		fmt.Println(minTemp)
	}
}
