package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		departures, employees int
		settings              string
	)

	const (
		lowerBoundary = 15
		upperBoundary = 30
	)

	reader := bufio.NewReader(os.Stdin)

	departuresTemp, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(-1)
		return
	}
	departuresTemp = strings.TrimSpace(departuresTemp)
	departures, err = strconv.Atoi(departuresTemp)
	if err != nil {
		fmt.Println(-1)
		return
	}

	for range departures {
		var employeesTemp string
		employeesTemp, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(-1)
			continue
		}
		employeesTemp = strings.TrimSpace(employeesTemp)
		employees, err = strconv.Atoi(employeesTemp)
		if err != nil {
			fmt.Println(-1)
			continue
		}

		lowerBorder, upperBorder := lowerBoundary, upperBoundary
		hasError := false

		for range employees {
			settings, err = reader.ReadString('\n')
			if err != nil {
				hasError = true
				continue
			}
			settings = strings.TrimSpace(settings)
			parts := strings.Fields(settings)
			sign := parts[0]
			temperature, _ := strconv.Atoi(parts[1])

			if temperature < lowerBorder || temperature > upperBorder {
				hasError = true
				continue
			}
			switch sign {
			case ">=":
				lowerBorder = max(lowerBorder, temperature)
			case "<=":
				upperBorder = min(upperBorder, temperature)
			default:
				hasError = true
				continue
			}

			if lowerBorder > upperBorder {
				hasError = true
				continue
			}
		}

		if hasError || lowerBorder > upperBorder {
			fmt.Println(-1)
		} else {
			fmt.Println(lowerBorder)
		}
	}
}


