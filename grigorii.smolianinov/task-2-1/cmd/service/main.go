package main

import (
	"fmt"
)

func main() {
	var N, K int
	fmt.Scan(&N)
	fmt.Scan(&K)

	for i := 0; i < N; i++ {
		minTemp := 15
		maxTemp := 30

		for j := 0; j < K; j++ {
			var sign string
			var t int
			fmt.Scan(&sign, &t)

			if sign == ">=" {
				if t > minTemp {
					minTemp = t
				}
			} else if sign == "<=" {
				if t < maxTemp {
					maxTemp = t
				}
			}
		}

		if minTemp <= maxTemp {
			fmt.Println(minTemp)
		} else {
			fmt.Println(-1)
		}
	}
}
