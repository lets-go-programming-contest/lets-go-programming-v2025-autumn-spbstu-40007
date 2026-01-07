package main

import (
	"fmt"
)

func main() {
	var officeCount int

	_, err := fmt.Scanln(&officeCount)
	if err != nil || officeCount < 1 || officeCount > 1000 {
		fmt.Println(-1)

		return
	}

	for i := range officeCount {
		var workersCount int

		_, err = fmt.Scanln(&workersCount)
		if err != nil || workersCount < 1 || workersCount > 1000 {
			fmt.Println(-1)

			return
		}

		lowerBound := 15
		upperBound := 30

		for j := range workersCount {
			var operator string
			var temp int

			_, err = fmt.Scanln(&operator, &temp)
			if err != nil || (operator != ">=" && operator != "<=") {
				fmt.Println(-1)

				return
			}

			if temp < 15 || temp > 30 {
				fmt.Println(-1)

				return
			}

			if operator == ">=" {
				if temp > lowerBound {
					lowerBound = temp
				}
			} else {
				if temp < upperBound {
					upperBound = temp
				}
			}

			if lowerBound > upperBound {
				fmt.Println(-1)
			} else {
				fmt.Println(lowerBound)
			}
		}
	}
}
