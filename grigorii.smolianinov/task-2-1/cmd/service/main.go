package main

import (
	"fmt"
	"os"
)

func main() {
	var departmentsCount, employeesCount int

	if _, err := fmt.Fscan(os.Stdin, &departmentsCount); err != nil {
		return
	}

	for range departmentsCount {
		if _, err := fmt.Fscan(os.Stdin, &employeesCount); err != nil {
			return
		}

		minTemp, maxTemp := 15, 30
		rangeValid := true

		for range employeesCount {
			var (
				sign string
				temp int
			)

			if _, err := fmt.Fscan(os.Stdin, &sign, &temp); err != nil {
				return
			}

			if !rangeValid {
				fmt.Println(-1)

				continue
			}

			if sign == ">=" && temp > minTemp {
				minTemp = temp
			}

			if sign == "<=" && temp < maxTemp {
				maxTemp = temp
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			}

			fmt.Println(-1)

			rangeValid = false
		}
	}
}
