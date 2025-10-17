package main

import (
	"fmt"
	"os"
)

func intMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func intMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func processDepartment(numRequests int) {
	var (
		lowBorder       = 15
		highBorder      = 30
		inContradiction bool
	)

	for range numRequests {
		var (
			temperature int
			sign        string
		)

		_, err := fmt.Scanf("%s %d", &sign, &temperature)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			return
		}

		if inContradiction {
			fmt.Println("-1")
			continue
		}

		if sign == ">=" {
			lowBorder = intMax(lowBorder, temperature)
		} else if sign == "<=" {
			highBorder = intMin(highBorder, temperature)
		}

		if lowBorder > highBorder {
			inContradiction = true
			fmt.Println("-1")
		} else {
			fmt.Println(lowBorder)
		}
	}
}

func main() {
	var (
		numDepartments int
		numRequests    int
	)

	_, err := fmt.Scan(&numDepartments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	for range numDepartments {
		_, err = fmt.Scan(&numRequests)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		processDepartment(numRequests)
	}
}
