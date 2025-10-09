package main

import (
	"fmt"
	"os"
)

func main() {
	var departmentsCount, employeesCount, temperatureValue int
	var sign string

	if _, err := fmt.Fscan(os.Stdin, &departmentsCount, &employeesCount); err != nil {
		return
	}

	for range departmentsCount {
		minTemperature := 15
		maxTemperature := 30

		for range employeesCount {

			if _, err := fmt.Fscan(os.Stdin, &sign, &temperatureValue); err != nil {
				return
			}

			if sign == ">=" {
				if temperatureValue > minTemperature {
					minTemperature = temperatureValue
				}
			} else if sign == "<=" {
				if temperatureValue < maxTemperature {
					maxTemperature = temperatureValue
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
