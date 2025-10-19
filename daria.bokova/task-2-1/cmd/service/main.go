package main

import "fmt"

func processEmployeeConstraint(operation string, temperature, minTemp, maxTemp int) (int, int, bool) {
	switch operation {
	case ">=":
		if temperature > maxTemp {
			return 0, 0, false
		}

		if temperature > minTemp {
			minTemp = temperature
		}

	case "<=":
		if temperature < minTemp {
			return 0, 0, false
		}

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
		valid := true

		for range employees {
			var (
				operation   string
				temperature int
			)

			if _, err := fmt.Scan(&operation, &temperature); err != nil {
				return
			}

			if !valid {
				fmt.Println(-1)

				continue
			}

			newMin, newMax, ok := processEmployeeConstraint(operation, temperature, minTemp, maxTemp)
			if !ok {
				fmt.Println(-1)

				valid = false

				continue
			}

			minTemp, maxTemp = newMin, newMax
			fmt.Println(minTemp)
		}
	}
}
