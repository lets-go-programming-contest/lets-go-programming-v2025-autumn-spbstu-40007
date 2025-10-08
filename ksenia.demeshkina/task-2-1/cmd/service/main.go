package main

import "fmt"

func temperature(workers int) {
	minTemp := 15
	maxTemp := 30

	for iterator := range workers {
		var operator string
		var value int
		
		_, err := fmt.Scan(&operator, &value)
		if err != nil {
			fmt.Println("Ошибка ввода")

			return
		}

		if operator == "<=" {
			if value < maxTemp {
				maxTemp = value
			}
		} else if operator == ">=" {
			if value > minTemp {
				minTemp = value
			}
		}

		if minTemp > maxTemp {
			fmt.Println(-1)

			for j := iterator + 1; j < workers; j++ {
				_, err := fmt.Scan(&operator, &value)
				
				if err != nil {
					fmt.Println("Ошибка ввода")

					return
				}
				fmt.Println(-1)
			}

			return
		} else {
			fmt.Println(minTemp)
		}
	}
}

func main() {
	var department, workers int
	_, err := fmt.Scan(&department)

	if err != nil {
		fmt.Println("Error")
	}

	if department < 1 || department > 1000 {
		fmt.Println("Количество отделов в диапазоне от 1 до 1000")

		return
	}

	for range department {
		_, err = fmt.Scan(&workers)

		if err != nil {
			fmt.Println("Error")
	}

		if workers < 1 || workers > 1000 {
			fmt.Println("Количество сотрудников от 1 до 1000")

			return
	}

	temperature(workers)

	}
}