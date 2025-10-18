package main

import (
	"fmt"
	"strconv"
)

func takeTemperature(K int) {
	minTemp := 15
	maxTemp := 30

	for range K {

		var (
			operation string
			value     string
		)

		_, err3 := fmt.Scanln(&operation, &value)
		if err3 != nil {
			fmt.Println("Invalid operation or value")
			return
		}

		tempInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Error invalid value")
		}

		switch operation {
		case "<=":
			maxTemp = min(maxTemp, tempInt)
		case ">=":
			minTemp = max(minTemp, tempInt)
		}
		if minTemp > maxTemp {
			fmt.Println(-1)
		}
		if minTemp <= maxTemp {
			fmt.Println(minTemp)
		}

	}
}

func main() {

	var (
		N, K int
	)

	_, err1 := fmt.Scanln(&N)
	if err1 != nil {
		fmt.Println("Invalid departments")
		return
	}

	for range N {
		_, err2 := fmt.Scanln(&K)
		if err2 != nil {
			fmt.Println("Invalid employees")
			return
		}

		takeTemperature(K)
	}

}
