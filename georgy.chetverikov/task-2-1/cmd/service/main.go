package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() { //nolint:funlen
	const (
		lowerBoundary  = 15
		upperBoundary  = 30
		settingsLength = 2
	)

	reader := bufio.NewReader(os.Stdin)
	departmentsTemp, _ := reader.ReadString('\n')
	departments, _ := strconv.Atoi(strings.TrimSpace(departmentsTemp))

	for range departments {
		employeesTemp, _ := reader.ReadString('\n')
		employees, _ := strconv.Atoi(strings.TrimSpace(employeesTemp))

		lowerBorder, upperBorder := lowerBoundary, upperBoundary
		hasError := false
		processedEmployees := 0

		for employeesNumber := range employees {
			processedEmployees = employeesNumber

			if hasError {
				fmt.Println(-1)
				reader.ReadString('\n') //nolint:errcheck

				continue
			}

			settingsTemp, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(-1)

				hasError = true

				continue
			}

			settings := strings.TrimSpace(settingsTemp)
			parts := strings.Fields(settings)
			if len(parts) < settingsLength {
				fmt.Println(-1)

				hasError = true

				continue
			}

			sign, temperature := func() (string, int) {
				temperature, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println(-1)

					hasError = true
				}

				return parts[0], temperature
			}()

			if sign == ">=" {
				lowerBorder = max(lowerBorder, temperature)
			} else if sign == "<=" {
				upperBorder = min(upperBorder, temperature)
			}

			if lowerBorder > upperBorder {
				fmt.Println(-1)
				hasError = true
			} else {
				fmt.Println(lowerBorder)
			}
		}

		if hasError {
			remaining := employees - (processedEmployees + 1)
			for range remaining {
				reader.ReadString('\n') //nolint:errcheck
				fmt.Println(-1)
			}
		}
	}
}
