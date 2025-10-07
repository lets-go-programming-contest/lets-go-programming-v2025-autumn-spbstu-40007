package main

import "fmt"

func main() {
	var departmentsCount int

	_, err := fmt.Scan(&departmentsCount)

	if err != nil {
		fmt.Println(err)
	}

	for range departmentsCount {
		maximumTemp := 30
		minimumTemp := 15

		var employeesCount int

		_, err = fmt.Scan(&employeesCount)

		if err != nil {
			fmt.Println(err)
		}

		for range employeesCount {
			var (
				operator  string
				tempValue int
			)

			_, err = fmt.Scan(&operator, &tempValue)

			if err != nil {
				fmt.Println(err)
			}

			if operator == "<=" && tempValue < maximumTemp {
				maximumTemp = tempValue
			}

			if operator == ">=" && tempValue > minimumTemp {
				minimumTemp = tempValue
			}

			if minimumTemp <= maximumTemp {
				fmt.Println(minimumTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
