package main

import (
	"fmt"
)

func main() {
	var (
		numberDepartments int
	)

	_, err := fmt.Scan(&numberDepartments)
	if err != nil {
		fmt.Println("Error reading input: ", err)

		return
	}

	for range numberDepartments {
		var (
			numberEmployees, temperature int
			operand                      string
		)

		_, err := fmt.Scan(&numberEmployees)
		if err != nil {
			fmt.Println("Error reading input: ", err)

			return
		}

		minTemperature, maxTemperature := 15, 30

		for range numberEmployees {
			_, err := fmt.Scan(&operand, &temperature)
			if err != nil {
				fmt.Println("Error reading input: ", err)

				return
			}

			if operand == "<=" {
				maxTemperature = min(maxTemperature, temperature)
			} else {
				minTemperature = max(temperature, minTemperature)
			}

			if minTemperature > maxTemperature {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemperature)
			}
		}
	}
}
