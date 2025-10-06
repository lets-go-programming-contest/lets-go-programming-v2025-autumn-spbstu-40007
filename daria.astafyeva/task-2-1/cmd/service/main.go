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

	for i := 0; i < empCount; i++ {
		var op string
		var temp int
		if _, err := fmt.Scanf("%s %d", &op, &temp); err != nil || (op != ">=" && op != "<=") {
			fmt.Println(-1)
			continue
		}

		if op == ">=" {
			minTemp = intMaximum(minTemp, temp)
		} else if op == "<=" {
			maxTemp = intMinimum(maxTemp, temp)
		}

		if minTemp <= maxTemp && minTemp >= 15 && minTemp <= 30 {
			fmt.Println(minTemp)
		} else {
			fmt.Println(-1)
		}
	}
}

func main() {
	var deptCount, empCount int
	if _, err := fmt.Scan(&deptCount); err != nil || deptCount < 1 || deptCount > 1000 {
		fmt.Println(-1)
		return
	}

	for i := 0; i < deptCount; i++ {
		if _, err := fmt.Scan(&empCount); err != nil || empCount < 1 || empCount > 1000 {
			fmt.Println(-1)
			continue
		}
		processDepartment(empCount)
	}
}
