package main

import (
	"fmt"
)

func main() {
	var departmentsCount int

	if _, err := fmt.Scan(&departmentsCount); err != nil {
		fmt.Println("Input error: failed to read количество отделов -", err)

		return
	}

	for range make([]struct{}, departmentsCount) {
		var staffCount int

		if _, err := fmt.Scan(&staffCount); err != nil {
			fmt.Println("Input error: failed to read staff count -", err)

			return
		}

		minT := 15
		maxT := 30

		for range make([]struct{}, staffCount) {
			var (
				operator string
				valueT   int
			)

			if _, err := fmt.Scan(&operator, &valueT); err != nil {
				fmt.Println("Input error: failed to read operator and value of temperature -", err)

				return
			}

			switch operator {
			case ">=":
				if valueT > minT {
					minT = valueT
				}
			case "<=":
				if valueT < maxT {
					maxT = valueT
				}
			}

			if minT > maxT {
				fmt.Println(-1)

				// продолжаем обработку следующих сотрудников
				continue
			}

			fmt.Println(minT)
		}
	}
}
