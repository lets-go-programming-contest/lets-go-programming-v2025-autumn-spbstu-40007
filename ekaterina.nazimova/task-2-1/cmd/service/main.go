package main

import (
	"fmt"
)

func main() {
	var departmentsCount int

	if _, err := fmt.Scan(&departmentsCount); err != nil {
		return
	}

	for range departmentsCount {
		var employeesCount int

		if _, err := fmt.Scan(&employeesCount); err != nil {
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
				return
			}

			if operator == ">=" && tempValue > minTemp {
				minTemp = tempValue
			}

			if operator == "<=" && tempValue < maxTemp {
				maxTemp = tempValue
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
