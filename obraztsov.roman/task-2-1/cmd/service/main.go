package main

import (
	"fmt"
	"strconv"
)

func takeTemperature(employees int) {
	minTemp := 15
	maxTemp := 30

	for range employees {
		var (
			operation string
			value     string
		)
		_, err3 := fmt.Scanln(&operation, &value)
		if err3 != nil {
			fmt.Println("Invalid operation or value")
			return
		}

		tempInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Error invalid value")
		}

		switch operation {
		case "<=":
			maxTemp = min(maxTemp, tempInt)
		case ">=":
			minTemp = max(minTemp, tempInt)
		}
		result := minTemp > maxTemp
		
		if result {
			fmt.Println(-1)
		} else {
			fmt.Println(minTemp)
		}
	}
}

func main() {
	var (
		departments, employees int
	)
	_, err1 := fmt.Scanln(&departments)
	if err1 != nil {
		fmt.Println("Invalid departments")
		return
	}

	for range departments {
		_, err2 := fmt.Scanln(&employees)
		if err2 != nil {
			fmt.Println("Invalid employees")
			return
		}
		takeTemperature(employees)
	}
}
