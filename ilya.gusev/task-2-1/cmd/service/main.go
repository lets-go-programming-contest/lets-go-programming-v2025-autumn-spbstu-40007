package main

import (
	"fmt"
)

func main() {
	var officeCount int

	_, err := fmt.Scanln(&officeCount)
	if err != nil || officeCount < 1 || officeCount > 1000 {
		fmt.Println(-1)

		return
	}

	for range officeCount {
		processOffice()
	}
}

func processOffice() {
	var workersCount int

	_, err := fmt.Scanln(&workersCount)
	if err != nil || workersCount < 1 || workersCount > 1000 {
		fmt.Println(-1)

		return
	}

	lowerBound := 15
	upperBound := 30

	for range workersCount {
		lowerBound, upperBound = processWorker(lowerBound, upperBound)
		if lowerBound == -1 {
			return
		}
	}
}

func processWorker(lowerBound, upperBound int) (int, int) {
	var operator string

	var temp int

	_, err := fmt.Scanln(&operator, &temp)
	if err != nil || (operator != ">=" && operator != "<=") || temp < 15 || temp > 30 {
		fmt.Println(-1)

		return -1, -1
	}

	if operator == ">=" && temp > lowerBound {
		lowerBound = temp
	}

	if operator == "<=" && temp < upperBound {
		upperBound = temp
	}

	if lowerBound > upperBound {
		fmt.Println(-1)
	} else {
		fmt.Println(lowerBound)
	}

	return lowerBound, upperBound
}
