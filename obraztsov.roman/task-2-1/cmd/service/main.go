package main

import (
	"fmt"
	"strconv"
)

func takeTemperature(k int) {
	minTemp := 15
	maxTemp := 30

	for range k {

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
		n, k int
	)

	_, err1 := fmt.Scanln(&n)
	if err1 != nil {
		fmt.Println("Invalid departments")
		return
	}

	for range n {
		_, err2 := fmt.Scanln(&k)
		if err2 != nil {
			fmt.Println("Invalid employees")
			return
		}

		takeTemperature(k)
	}

}
