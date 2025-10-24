package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	minTemp       = 15
	maxTemp       = 30
	expectedParts = 2
)

func processDepartment(scanner *bufio.Scanner, numEmployees int) {
	currentMin := minTemp
	currentMax := maxTemp

	for range make([]struct{}, numEmployees) {
		if !scanner.Scan() {
			return
		}

		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, " ")
		if len(parts) != expectedParts {
			fmt.Println(-1)
			continue
		}

		operator := parts[0]
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println(-1)
			continue
		}

		switch operator {
		case ">=":
			if value > currentMin {
				currentMin = value
			}
		case "<=":
			if value < currentMax {
				currentMax = value
			}
		default:
			fmt.Println(-1)
			continue
		}

		if currentMin <= currentMax {
			fmt.Println(currentMin)
		} else {
			fmt.Println(-1)
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	if !scanner.Scan() {
		return
	}

	numDepartments, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || numDepartments < 1 || numDepartments > 1000 {
		return
	}

	for range make([]struct{}, numDepartments) {
		if !scanner.Scan() {
			return
		}

		numEmployees, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || numEmployees < 1 || numEmployees > 1000 {
			return
		}

		processDepartment(scanner, numEmployees)
	}
}
