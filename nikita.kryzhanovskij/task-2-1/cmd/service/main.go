package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseBorder(border string) (string, int, error) {
	fields := strings.Fields(border)

	if fields[0] != "<=" && fields[0] != ">=" {
		return "", 0, fmt.Errorf("invalid operator")
	}

	temp, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid temperature")
	}

	if temp > 30 || temp < 15 {
		return "", 0, fmt.Errorf("invalid temperature")
	}

	return fields[0], temp, nil
}

func main() {
	var n int

	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		var k int
		fmt.Scan(&k)

		minTemp := 15
		maxTemp := 30

		for j := 0; j < k; j++ {
			var (
				op   string
				temp string
			)

			fmt.Scan(&op, &temp)

			tempBorder := fmt.Sprintf("%s %s", op, temp)

			operation, temperature, err := parseBorder(tempBorder)

			if err != nil {
				fmt.Println("Error has occurred:", err)
				return
			}

			switch operation {
			case ">=":
				minTemp = max(minTemp, temperature)
			case "<=":
				maxTemp = min(maxTemp, temperature)
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
				break
			}

			fmt.Println(minTemp)
		}
	}
}
