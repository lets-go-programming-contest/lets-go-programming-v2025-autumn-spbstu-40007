package main

import "fmt"

func temperature(workers int) error {
	minTemp := 15
	maxTemp := 30

	for iterator := range workers {
		var operator string

		var value int

		_, err := fmt.Scan(&operator, &value)
		if err != nil {
			return err
		}

		if operator == "<=" {
			if value < maxTemp {
				maxTemp = value
			}
		}

		if operator == ">=" {
			if value > minTemp {
				minTemp = value
			}
		}

		if minTemp > maxTemp {
			fmt.Println(-1)

			for j := iterator + 1; j < workers; j++ {
				_, err := fmt.Scan(&operator, &value)
				if err != nil {
					return err
				}

				fmt.Println(-1)
			}

			return nil
		}

		fmt.Println(minTemp)
	}

	return nil
}

func main() {
	var department, workers int

	_, err := fmt.Scan(&department)
	if err != nil {
		fmt.Println("Invalid department input format")
	}

	if department < 1 || department > 1000 {
		fmt.Println("The number of departments must be in the range from 1 to 1000")

		return
	}

	for range department {
		_, err = fmt.Scan(&workers)
		if err != nil {
			fmt.Println("Invalid employee input format")

			return
		}

		if workers < 1 || workers > 1000 {
			fmt.Println("The number of employees must be from 1 to 1000")

			return
		}

		err := temperature(workers)
		if err != nil {
			fmt.Println("temperature data entry error")

			return
		}
	}
}
