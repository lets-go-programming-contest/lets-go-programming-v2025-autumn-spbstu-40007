package main

import "fmt"

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
		isValid := true

		for range employees {
			var operation string
			var temperature int

			if _, err := fmt.Scan(&operation, &temperature); err != nil {
				return
			}

			if !isValid {
				fmt.Println(-1)

				continue
			}

			switch operation {
			case ">=":
				if temperature > maxTemp {
					isValid = false
				} else if temperature > minTemp {
					minTemp = temperature
				}
			case "<=":
				if temperature < minTemp {
					isValid = false
				} else if temperature < maxTemp {
					maxTemp = temperature
				}
			}

			if isValid {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
