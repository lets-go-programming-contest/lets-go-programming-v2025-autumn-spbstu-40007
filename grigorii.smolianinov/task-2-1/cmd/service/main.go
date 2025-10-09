package main

import (
	"fmt"
	"os"
)

func main() {
	var departmentsCount, employeesCount int

	minTemp, maxTemp := 15, 30
	rangeValid := true

	if _, err := fmt.Fscan(os.Stdin, &departmentsCount, &employeesCount); err != nil {
		return
	}

	for range departmentsCount {
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
			} else {
				fmt.Println(-1)

				rangeValid = false
			}
		}
	}
}
