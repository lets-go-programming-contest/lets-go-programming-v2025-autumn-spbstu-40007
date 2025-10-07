package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	
	departuresTemp, _ := reader.ReadString('\n')
	departuresTemp = strings.TrimSpace(departuresTemp)
	departures, _ := strconv.Atoi(departuresTemp)

	for i := 0; i < departures; i++ {
		employeesTemp, _ := reader.ReadString('\n')
		employeesTemp = strings.TrimSpace(employeesTemp)
		employees, _ := strconv.Atoi(employeesTemp)

		lowerBorder, upperBorder := 15, 30

		for j := 0; j < employees; j++ {
			settings, _ := reader.ReadString('\n')
			settings = strings.TrimSpace(settings)
			
			parts := strings.Fields(settings)
			if len(parts) < 2 {
				if lowerBorder <= upperBorder {
					fmt.Println(lowerBorder)
				} else {
					fmt.Println(-1)
				}
				continue
			}
			
			sign := parts[0]
			number, _ := strconv.Atoi(parts[1])

			switch sign {
			case ">=":
				if number > lowerBorder {
					lowerBorder = number
				}
			case "<=":
				if number < upperBorder {
					upperBorder = number
				}
			default:
				if lowerBorder <= upperBorder {
					fmt.Println(lowerBorder)
				} else {
					fmt.Println(-1)
				}
				continue
			}

			if lowerBorder > upperBorder {
				fmt.Println(-1)
				for k := j + 1; k < employees; k++ {
					reader.ReadString('\n')
				}
				break
			}
			
			fmt.Println(lowerBorder)
		}
	}
}
