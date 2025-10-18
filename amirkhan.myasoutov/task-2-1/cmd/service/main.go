package main

import "fmt"

func main() {
	var departmentAmount int

	_, err := fmt.Scan(&departmentAmount)
	if err != nil {
		return
	}

	for range departmentAmount {
		maximumTemp := 30
		minimumTemp := 15

		var employeeAmount int

		_, err = fmt.Scan(&employeeAmount)
		if err != nil {
			return
		}

		for range employeeAmount {
			var (
				operand     string
				temperature int
			)

			_, err = fmt.Scan(&operand, &temperature)
			if err != nil {
				return
			}

			if operand == "<=" && temperature < maximumTemp {
				maximumTemp = temperature
			}

			if operand == ">=" && temperature > minimumTemp {
				minimumTemp = temperature
			}

			if minimumTemp <= maximumTemp {
				fmt.Println(minimumTemp)

				continue
			}

			fmt.Println(-1)
		}
	}
}
