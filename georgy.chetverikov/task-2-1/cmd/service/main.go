package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	const (
		lowerBoundary = 15
		upperBoundary = 30
	)

	var departments int

	reader := bufio.NewReader(os.Stdin)
	departmentsTemp, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(-1)
		return
	}

	departments, err = strconv.Atoi(strings.TrimSpace(departmentsTemp))
	if err != nil {
		fmt.Println(-1)
		return
	}

	for departmentIndex := 0; departmentIndex < departments; departmentIndex++ {
		var employees int

		employeesTemp, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(-1)
			return
		}

		employees, err = strconv.Atoi(strings.TrimSpace(employeesTemp))
		if err != nil {
			fmt.Println(-1)
			return
		}
		
		lowerBorder, upperBorder := lowerBoundary, upperBoundary
		hasError := false
		currentEmployee := 0

		for employeeIndex := 0; employeeIndex <= employees; employeeIndex++ {
			currentEmployee = employeeIndex
			var settings string

			settingsTemp, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(-1)
				return
			}

			settings = strings.TrimSpace(settingsTemp)

			sign, temperature := func() (string, int) {
				parts := strings.Fields(settings)
				if len(parts) < 2 {
					return "", 0
				}
				temperature, err := strconv.Atoi(parts[1])
				if err != nil {
					return "", 0
				}
				return parts[0], temperature
			}()

			if sign == "" {
				fmt.Println(-1)
				hasError = true
				break
			}

			if sign == ">=" {
				lowerBorder = max(lowerBorder, temperature)
			} else if sign == "<=" {
				upperBorder = min(upperBorder, temperature)
			}

			if lowerBorder > upperBorder {
				fmt.Println(-1)
				hasError = true
				break
			} else {
				fmt.Println(lowerBorder)
			}
		}
		
		if hasError {
			remainingEmployees := employees - (currentEmployee + 1)
			for i := 0; i < remainingEmployees; i++ {
				reader.ReadString('\n')
			}
		}
	}
}

