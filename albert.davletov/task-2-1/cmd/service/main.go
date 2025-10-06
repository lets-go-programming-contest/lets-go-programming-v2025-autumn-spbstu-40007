package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readInput(scanner *bufio.Scanner) string {
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		fmt.Println("Error reading input")
		os.Exit(0)
	}

	return scanner.Text()
}

func main() {
	minTemperature := 15
	maxTemperature := 30

	reader := bufio.NewScanner(os.Stdin)

	numberDepartments, _ := strconv.Atoi(readInput(reader))
	for range numberDepartments {
		numberEmployees, _ := strconv.Atoi(readInput(reader))
		for range numberEmployees {
			preferences := readInput(reader)
			operand := preferences[0]
			temperature, _ := strconv.Atoi(preferences[3:])

			if operand == '<' {
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
