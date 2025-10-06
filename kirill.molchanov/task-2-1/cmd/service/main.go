package main

import (
	"fmt"
)

func getTemperature(k int) { //nolint:varnamelen
	var (
		temperature int
		sign        string
		lowBorder   = 30
		highBorder  = 15
	)

	for range k {
		_, err := fmt.Scanln(&sign, &temperature)
		if err != nil {
			fmt.Println(-1)
		}

		if !(15 <= temperature && temperature <= 30) {
			return
		}

		switch sign {
		case ">=":
			if temperature > highBorder {
				highBorder = temperature
			}
		case "<=":
			if temperature < lowBorder {
				lowBorder = temperature
			}
		default:
			return
		}

		if highBorder <= lowBorder {
			fmt.Println(highBorder)
		} else {
			fmt.Println(-1)
		}
	}
}

func main() {
	var (
		n, k int //nolint:varnamelen
	)

	_, err := fmt.Scan(&n)
	if err != nil {
		return
	}

	for range n {
		_, err := fmt.Scan(&k)
		if err != nil {
			return
		}

		getTemperature(k)
	}
}
