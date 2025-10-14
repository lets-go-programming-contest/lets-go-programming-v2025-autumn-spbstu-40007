package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//nolint:gci
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
			settingsTemp, err := reader.ReadString('\n') //nolint:wsl
			if err != nil {                              //nolint:wsl
				fmt.Println(-1)
				hasError = true //nolint:wsl

				continue
			}
			settings := strings.TrimSpace(settingsTemp) //nolint:wsl
			parts := strings.Fields(settings)
			if len(parts) < settingsLength { //nolint:wsl
				fmt.Println(-1)
				hasError = true //nolint:wsl

				continue
			}
			sign, temperature := func() (string, int) { //nolint:wsl
				temperature, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println(-1)
					hasError = true //nolint:wsl
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
				hasError = true //nolint:wsl
			} else {
				fmt.Println(lowerBorder)
			}
		}
		if hasError { //nolint:wsl
			remaining := employees - (processedEmployees + 1)
			for range remaining {
				reader.ReadString('\n') //nolint:errcheck
				fmt.Println(-1)
			}
		}
	}
}
