package main

import (
	"fmt"
)

func main() {
	var departmentsCount int

	if _, err := fmt.Scan(&departmentsCount); err != nil {
		fmt.Println("Input error: failed to read department count -", err)
		return
	}

	for range departmentsCount {
		var employeesCount int

		if _, err := fmt.Scan(&employeesCount); err != nil {
			fmt.Println("Input error: failed to read employees count -", err)
			return
		}

		minTemp := 15
		maxTemp := 30

		for range employeesCount {
			var (
				operator  string
				tempValue int
			)

			if _, err := fmt.Scan(&operator, &tempValue); err != nil {
				fmt.Println("Input error: failed to read operator and temperature value -", err)
				return
			}

			if operator == ">=" && tempValue > minTemp {
				minTemp = tempValue
			}

			if operator == "<=" && tempValue < maxTemp {
				maxTemp = tempValue
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
				return
			}
			fmt.Println(minTemp)

		}
	}
}
