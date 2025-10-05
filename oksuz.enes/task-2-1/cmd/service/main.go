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
		var temperature int
		if _, err := fmt.Fscan(reader, &operator, &temperature); err != nil {
			return
		}
		if operator == ">=" {
			if temperature > minTemp {
				minTemp = temperature
			}
		} else if operator == "<=" {
			if temperature < maxTemp {
				maxTemp = temperature
			}
		}
		if minTemp > maxTemp {
			fmt.Println(-1)
		} else {
			fmt.Println(minTemp)
		}
	}
}
