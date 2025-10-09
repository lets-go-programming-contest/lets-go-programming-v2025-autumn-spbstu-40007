package main

import (
	"fmt"
	"os"
)

func main() {
	var departmentsCount, employeesCount int

	if _, err := fmt.Fscan(os.Stdin, &departmentsCount, &employeesCount); err != nil {
		return
	}

	for departmentIndex := 0; departmentIndex < departmentsCount; departmentIndex++ {
		minTemperature := 15
		maxTemperature := 30

		for employeeIndex := 0; employeeIndex < employeesCount; employeeIndex++ {
			var sign string
			var temperature int

			if _, err := fmt.Fscan(os.Stdin, &sign, &temperature); err != nil {
				return
			}

			if sign == ">=" {
				if temperature > minTemperature {
					minTemperature = temperature
				}
			} else if sign == "<=" {
				if temperature < maxTemperature {
					maxTemperature = temperature
				}
			}

			if minTemperature <= maxTemperature {
				fmt.Println(minTemperature)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
