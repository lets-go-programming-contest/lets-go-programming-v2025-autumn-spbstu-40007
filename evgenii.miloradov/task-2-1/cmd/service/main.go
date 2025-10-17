package main

import (
	"fmt"
)

func main() {
	var numOfState, numOfEmp int
	maxTemp := 15
	minTemp := 30

	fmt.Scan(&numOfState)

	for i := 0; i < numOfState; i++ {
		maxTemp = 15
		minTemp = 30

		fmt.Scan(&numOfEmp)

		for j := 0; j < numOfEmp; j++ {

			var operand string
			var comfTemp int

			_, err := fmt.Scan(&operand, &comfTemp)
			if err != nil {
				fmt.Println("operand")
				return
			}

			switch operand {
			case ">=":
				if comfTemp > maxTemp {
					maxTemp = comfTemp
				}
			case "<=":
				if comfTemp < minTemp {
					minTemp = comfTemp
				}
			default:
				fmt.Println("-1")

			}

			if maxTemp <= minTemp {
				fmt.Println(maxTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
