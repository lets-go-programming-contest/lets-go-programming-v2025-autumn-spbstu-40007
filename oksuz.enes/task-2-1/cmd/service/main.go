package main

import (
	"fmt"
)

func main() {
	var room, workers, minTemp, maxTemp int

	if _, err := fmt.Scanln(&room, &workers); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	for i := 0; i < workers; i++ {
		var (
			operator string
			temp     int
		)

		if _, err := fmt.Scanln(&operator); err != nil {
			fmt.Println("Error reading operator:", err)
			return
		}

		if _, err := fmt.Scanln(&temp); err != nil {
			fmt.Println("Error reading temp:", err)
			return
		}
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
