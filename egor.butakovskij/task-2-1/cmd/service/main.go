package main

import (
	"fmt"
	"log"
)

func getRecommendedTemperature(employeesCount int) error {
	var (
		temp, recTemp, highBorder, lowBorder int
		sign                                 string
	)

	highBorder = 30
	lowBorder = 15

	for range employeesCount {
		_, err := fmt.Scanf("%s %d", &sign, &temp)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		if sign == ">=" {
			lowBorder = max(lowBorder, temp)
		}

		if sign == "<=" {
			highBorder = min(highBorder, temp)
		}

		if lowBorder > highBorder {
			recTemp = -1
		} else {
			recTemp = lowBorder
		}

		fmt.Println(recTemp)
	}

	return nil
}

func main() {
	log.SetFlags(0)

	var departmentsCount, employeesCount int

	_, err := fmt.Scan(&departmentsCount)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for range departmentsCount {
		_, err := fmt.Scan(&employeesCount)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = getRecommendedTemperature(employeesCount)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}
