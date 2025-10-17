package main

import "fmt"

func processEmployeeConstraint(operation string, temperature, minTemp, maxTemp int) (int, int, bool) {
	if temperature < minTemp || temperature > maxTemp {
		return minTemp, maxTemp, false
	}

	switch operation {
	case ">=":
		if temperature > minTemp {
			minTemp = temperature
		}
	case "<=":
		if temperature < maxTemp {
			maxTemp = temperature
		}
	}

	return minTemp, maxTemp, true
}

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		return
	}

	for range departments {
		var employees int
		if _, err := fmt.Scan(&employees); err != nil {
			return
		}

		minTemp := 15
		maxTemp := 30

		for range employees {
			var (
				operation   string
				temperature int
			)

			if _, err := fmt.Scan(&operation, &temperature); err != nil {
				return
			}

			newMin, newMax, isValid := processEmployeeConstraint(operation, temperature, minTemp, maxTemp)

			if !isValid {
				fmt.Println(-1)
				minTemp, maxTemp = newMin, newMax
				continue
			}

			minTemp, maxTemp = newMin, newMax
			fmt.Println(minTemp)
		}
	}
}
