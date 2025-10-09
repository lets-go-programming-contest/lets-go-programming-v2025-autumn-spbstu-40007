package main

import (
	"fmt"
	"os"
)

func main() {
	var (
		departmentsCount, employeesCount, temperatureValue int
		sign                                               string
	)

	if _, err := fmt.Fscan(os.Stdin, &departmentsCount, &employeesCount); err != nil {
		return
	}

	for i := 0; i < departmentsCount; i++ {
		minTemperature := 15
		maxTemperature := 30
		valid := true

		for j := 0; j < employeesCount; j++ {
			if _, err := fmt.Fscan(os.Stdin, &sign, &temperatureValue); err != nil {
				return
			}

			if valid {
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
					valid = false
				}
			} else {
				fmt.Println(-1)
			}
		}
	}
}
