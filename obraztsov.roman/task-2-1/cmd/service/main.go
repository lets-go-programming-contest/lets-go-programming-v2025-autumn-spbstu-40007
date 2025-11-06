package main

import (
	"fmt"
	"strconv"
)

func takeTemperature(employees int) {
	minTemp := 15
	maxTemp := 30

	for i := 0; i < employees; i++ {
		var (
			operation string
			value     string
		)

		_, err := fmt.Scanln(&operation, &value)
		if err != nil {
			fmt.Println("Invalid operation or value")
			return
		}

		tempInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Error invalid value")
			return
		}

		switch operation {
		case "<=":
			if tempInt < maxTemp {
				maxTemp = tempInt
			}
		case ">=":
			if tempInt > minTemp {
				minTemp = tempInt
			}
		}

		if minTemp > maxTemp {
			fmt.Println(-1)
		} else {
			fmt.Println(minTemp)
		}
	}
}

func main() {
	var departments, employees int

	_, err := fmt.Scanln(&departments)
	if err != nil {
		fmt.Println("Invalid departments")
		return
	}

	for i := 0; i < departments; i++ {
		_, err := fmt.Scanln(&employees)
		if err != nil {
			fmt.Println("Invalid employees")
			return
		}

		takeTemperature(employees)
	}
}
