package main

import (
	"fmt"
)

func getTemperature(k int) {
	var (
		temperature int
		sign        string
		lowBorder   = 30
		highBorder  = 15
	)

	for range k {
		fmt.Scanln(&sign, &temperature)

		if !(15 <= temperature && temperature <= 30) {
			fmt.Println(-1)
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
			fmt.Println(-1)
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
		n, k int
	)

	fmt.Scanln(&n)

	if !(1 <= n && n <= 1000) {
		fmt.Println(-1)
	}

	for range n {
		fmt.Scanln(&k)

		if !(1 <= k && k <= 1000) {
			fmt.Println(-1)
		}

		getTemperature(k)
	}
}
