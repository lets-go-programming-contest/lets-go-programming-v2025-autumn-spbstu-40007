package main

import (
	"fmt"
)

func main() {
	var numOfState, numOfEmp int

	var maxTemp, minTemp int

	if _, err := fmt.Scan(&numOfState); err != nil {
		fmt.Println(-1)

		return
	}
	//nolint:(intrange)
	for i := 0; i < numOfState; i++ {

		maxTemp = 15
		minTemp = 30

		if _, err := fmt.Scan(&numOfEmp); err != nil {
			fmt.Println(-1)

			return
		}
		//nolint:(intrange)
		for j := 0; j < numOfEmp; j++ {
			var operand string

			var comfTemp int

			_, err := fmt.Scan(&operand, &comfTemp)
			if err != nil {
				fmt.Println("-1")

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

				return
			}

			if maxTemp <= minTemp {
				fmt.Println(maxTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
