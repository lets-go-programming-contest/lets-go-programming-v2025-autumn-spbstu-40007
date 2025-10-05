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
	var minTemperature int = 15
	var maxTemperature int = 30

	reader := bufio.NewScanner(os.Stdin)
	numberDepartments, _ := strconv.Atoi(readInput(reader))
	for i := 0; i < numberDepartments; i++ {
		numberEmployees, _ := strconv.Atoi(readInput(reader))
		for v := 0; v < numberEmployees; v++ {
			preferences := readInput(reader)
			operand := preferences[0]
			temperature, _ := strconv.Atoi(preferences[3:])
			if operand == '<' {
				maxTemperature = min(maxTemperature, temperature)
			} else {
				if temperature > maxTemperature {
					fmt.Println(-1)
					break
				}
				minTemperature = max(temperature, minTemperature)
			}
			fmt.Println(minTemperature)
		}
		minTemperature = 15
		maxTemperature = 30
	}
}
