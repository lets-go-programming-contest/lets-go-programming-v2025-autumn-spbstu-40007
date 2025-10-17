package main

import (
	"fmt"
)

func main() { //nolint:cyclop
	var departmentNum, employeesNum int

	_, err := fmt.Scanln(&departmentNum)
	if err != nil || departmentNum < 1 || departmentNum > 1000 {
		fmt.Println(-1)

		return
	}

	for range departmentNum {
		_, err = fmt.Scanln(&employeesNum)
		if err != nil || employeesNum < 1 || employeesNum > 1000 {
			fmt.Println(-1)

			return
		}

		minTemperature := 15
		maxTemperature := 30

		for range employeesNum {
			var (
				sign           string
				curTemperature int
			)

			_, err = fmt.Scanln(&sign, &curTemperature)
			if err != nil || (sign != ">=" && sign != "<=") || curTemperature < 15 || curTemperature > 30 {
				fmt.Println(-1)

				return
			}

			if sign == ">=" && curTemperature > minTemperature {
				minTemperature = curTemperature
			}

			if sign == "<=" && curTemperature < maxTemperature {
				maxTemperature = curTemperature
			}

			if maxTemperature < minTemperature {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemperature)
			}
		}
	}
}
