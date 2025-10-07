package main

import "fmt"

func main() {
	var departmentAmount int
	fmt.Scan(&departmentAmount)

	for i := 0; i < departmentAmount; i++ {
		maximumTemp := 30
		minimumTemp := 15

		var employeeAmount int
		fmt.Scan(&employeeAmount)

		for j := 0; j < employeeAmount; j++ {
			var (
				operand     string
				temperature int
			)

			fmt.Scan(&operand, &temperature)

			if operand == "<=" && temperature < maximumTemp {
				maximumTemp = temperature
			}
			if operand == ">=" && temperature > minimumTemp {
				minimumTemp = temperature
			}

			if minimumTemp <= maximumTemp {
				fmt.Println(minimumTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}

}
