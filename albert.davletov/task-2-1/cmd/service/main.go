package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readInput(scanner *bufio.Scanner) (string, error) {
	scanner.Scan()
	err := scanner.Err()

	return scanner.Text(), err
}

func main() {
	minTemperature := 15
	maxTemperature := 30

	reader := bufio.NewScanner(os.Stdin)

	numberDepartmentsString, err := readInput(reader)
	if err != nil {
		fmt.Println(err)

		return
	}

	numberDepartments, err := strconv.Atoi(numberDepartmentsString)
	if err != nil {
		fmt.Println(err)

		return
	}

	for range numberDepartments {
		numberEmployeesString, err := readInput(reader)
		if err != nil {
			fmt.Println(err)

			return
		}

		numberEmployees, err := strconv.Atoi(numberEmployeesString)
		if err != nil {
			fmt.Println(err)

			return
		}

		for range numberEmployees {
			preferences, err := readInput(reader)
			if err != nil {
				fmt.Println(err)

				return
			}

			operand := preferences[:2]
			temperature, err := strconv.Atoi(preferences[3:])
			if err != nil {
				fmt.Println(err)

				return
			}

			if operand == "<=" {
				maxTemperature = min(maxTemperature, temperature)
			} else {
				minTemperature = max(temperature, minTemperature)
			}

			if minTemperature > maxTemperature {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemperature)
			}
		}

		minTemperature = 15
		maxTemperature = 30
	}
}
