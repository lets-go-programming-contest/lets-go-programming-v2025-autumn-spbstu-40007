package main

import (
	"fmt"
)

func findOptimalTemp(numOfEmployees int) {
	minTemperature := 15
	maxTemperature := 30

	for range numOfEmployees {
		var (
			sign        string
			temperature int
		)

		_, err := fmt.Scanln(&sign, &temperature)
		if err != nil || (sign != ">=" && sign != "<=") {
			fmt.Println(-1)

			continue
		}

		if sign == ">=" {
			minTemperature = max(temperature, minTemperature)
		}

		if sign == "<=" {
			maxTemperature = min(temperature, maxTemperature)
		}

		if minTemperature > maxTemperature {
			fmt.Println(-1)
		} else {
			fmt.Println(minTemperature)
		}
	}
}

func main() {
	var numOfDepartments, numOfEmployees int

	_, err := fmt.Scanln(&numOfDepartments)
	if err != nil || numOfDepartments > 1000 || numOfDepartments < 1 {
		fmt.Println(-1)

		return
	}

	for range numOfDepartments {
		_, err = fmt.Scanln(&numOfEmployees)
		if err != nil || numOfEmployees > 1000 || numOfEmployees < 1 {
			fmt.Println(-1)

			continue
		}

		findOptimalTemp(numOfEmployees)
	}
}
