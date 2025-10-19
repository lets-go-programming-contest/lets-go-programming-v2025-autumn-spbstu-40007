package main

import "fmt"

func intMaximum(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func intMinimum(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func processDepartment(empCount int) {
	minTemp, maxTemp := 15, 30

	for range empCount {
		var (
			operator string
			temp     int
		)

		if _, err := fmt.Scanf("%s %d", &operator, &temp); err != nil || (operator != ">=" && operator != "<=") {
			fmt.Println(-1)

			continue
		}

		if operator == ">=" {
			minTemp = intMaximum(minTemp, temp)
		}

		if operator == "<=" {
			maxTemp = intMinimum(maxTemp, temp)
		}

		if minTemp <= maxTemp && minTemp >= 15 && minTemp <= 30 {
			fmt.Println(minTemp)

			continue
		}

		fmt.Println(-1)
	}
}

func main() {
	var deptCount, empCount int
	if _, err := fmt.Scan(&deptCount); err != nil || deptCount < 1 || deptCount > 1000 {
		fmt.Println(-1)

		return
	}

	for range deptCount {
		if _, err := fmt.Scan(&empCount); err != nil || empCount < 1 || empCount > 1000 {
			fmt.Println(-1)

			continue
		}

		processDepartment(empCount)
	}
}
