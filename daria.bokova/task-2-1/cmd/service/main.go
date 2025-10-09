package main

import "fmt"

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		return
	}

	for i := 0; i < departments; i++ {
		var employees int
		if _, err := fmt.Scan(&employees); err != nil {
			return
		}

		minTemp := 15
		maxTemp := 30
		isValid := true

		for j := 0; j < employees; j++ {
			var operation string
			var temperature int
			if _, err := fmt.Scan(&operation, &temperature); err != nil {
				return
			}

			if !isValid {
				continue
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

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
				isValid = false
			}
		}
	}
}
