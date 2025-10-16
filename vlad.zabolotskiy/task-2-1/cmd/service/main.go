package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	minTempConst        = 15
	maxTempConst        = 30
	opMoreConst         = ">="
	opLessConst         = "<="
	expectedFieldsCount = 2
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	quantityDepartments, _ := reader.ReadString('\n')
	departments, _ := strconv.Atoi(strings.TrimSpace(quantityDepartments))

	if departments < 1 || departments > 1000 {
		fmt.Println("Invalid departments range")
		return
	}

	quantityEmployees, _ := reader.ReadString('\n')
	employees, _ := strconv.Atoi(strings.TrimSpace(quantityEmployees))

	optimalTemperature(reader, departments, employees)
}

func optimalTemperature(reader *bufio.Reader, departments, employees int) {
	for range departments {
		if employees < 1 || employees > 1000 {
			fmt.Println("Invalid employees range")
			return
		} else {
			minTemp := minTempConst
			maxTemp := maxTempConst
			opMore := opMoreConst
			opLess := opLessConst

			for range employees {
				preference, _ := reader.ReadString('\n')
				preference = strings.TrimSpace(preference)

				data := strings.Fields(preference)

				if len(data) != expectedFieldsCount {
					fmt.Println(-1)
				} else {
					operator := data[0]
					temperature, err := strconv.Atoi(data[1])

					if err != nil {
						fmt.Println(-1)
						continue
					}

					switch operator {
					case opMore:
						minTemp = max(minTemp, temperature)
					case opLess:
						maxTemp = min(maxTemp, temperature)
					default:
						fmt.Println(-1)
						continue
					}

					if minTemp > maxTemp {
						fmt.Println(-1)
					} else {
						fmt.Println(minTemp)
					}
				}
			}
		}
	}
}
