package main

import "fmt"

func getRecommendedTemperature(employeesCount int) {
	var (
		temp, recTemp, highBorder, lowBorder int
		sign                                 string
	)
	highBorder = 30
	lowBorder = 15

	for range employeesCount {
		_, err := fmt.Scanf("%s %d", &sign, &temp)
		if err != nil {
			fmt.Println(err)
		}

		if sign == ">=" {
			lowBorder = max(lowBorder, temp)
		}

		if sign == "<=" {
			highBorder = min(highBorder, temp)
		}

		if lowBorder > highBorder {
			recTemp = -1
		} else {
			recTemp = lowBorder
		}

		fmt.Println(recTemp)
	}
}

func main() {
	var (
		departmentsCount, employeesCount int
	)

	_, err := fmt.Scan(&departmentsCount)
	if err != nil {
		fmt.Println(err)
	}

	for range departmentsCount {
		_, err := fmt.Scan(&employeesCount)
		if err != nil {
			fmt.Println(err)
		}
		getRecommendedTemperature(employeesCount)
	}
}
