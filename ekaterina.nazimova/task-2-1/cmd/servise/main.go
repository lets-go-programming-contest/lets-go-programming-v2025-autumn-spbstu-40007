package main

import (
	"fmt"
)

func main() {
	var n int
	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		var k int
		fmt.Scan(&k)

		minTemp := 15
		maxTemp := 30

		for j := 0; j < k; j++ {
			var op string
			var x int
			fmt.Scan(&op, &x)

			if op == ">=" && x > minTemp {
				minTemp = x
			}
			if op == "<=" && x < maxTemp {
				maxTemp = x
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
