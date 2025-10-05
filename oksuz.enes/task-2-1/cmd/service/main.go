package main

import (
	"fmt"
)

func main() {
	var room, workers, minTemp, maxTemp int

	fmt.Scanln(&room, &workers)

	for i := 0; i < workers; i++ {
		var (
			operator string
			temp     int
		)

		fmt.Scanln(&operator)
		fmt.Scanln(&temp)

		if operator == ">=" {
			minTemp = max(minTemp, temp)

		} else if operator == "<=" {
			maxTemp = min(maxTemp, temp)
		}

	}

	if minTemp > maxTemp {
		fmt.Println(-1)
	} else {
		fmt.Println(minTemp)
	}

}
