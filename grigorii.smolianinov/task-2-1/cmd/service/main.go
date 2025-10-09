package main

import (
	"fmt"
	"os"
)

func main() {
	var N, K int

	if _, err := fmt.Fscan(os.Stdin, &N, &K); err != nil {
		return
	}

	for i := 0; i < N; i++ {
		minTemp := 15
		maxTemp := 30

		for j := 0; j < K; j++ {
			var sign string
			var t int

			if _, err := fmt.Fscan(os.Stdin, &sign, &t); err != nil {
				return
			}

			if sign == ">=" {
				if t > minTemp {
					minTemp = t
				}
			} else if sign == "<=" {
				if t < maxTemp {
					maxTemp = t
				}
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemp)
			}
		}
		fmt.Println()
	}
}
