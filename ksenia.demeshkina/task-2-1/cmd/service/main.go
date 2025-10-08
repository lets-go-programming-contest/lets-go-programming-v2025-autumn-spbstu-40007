package main

import "fmt"

func temperature(k int) {
	minTemp := 15
	maxTemp := 30

	for i := range k {
		var operator string;
		var value int;
		_, err := fmt.Scan(&operator, &value)
		if err != nil {
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

			for j := i + 1; j < k; j++ {
				fmt.Scan(&operator, &value)
				fmt.Println(-1)
			}

			return
		} else {
			fmt.Println(minTemp)
		}
	}
}

func main() {
	var N, k int
	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Println("Error")
	}

	if N < 1 || N > 1000 {
		fmt.Println("Количество отделов в диапазоне от 1 до 1000")

		return
	}

	for range N {
		_, err = fmt.Scan(&k)
		if err != nil {
			fmt.Println("Error")
	}

		if k < 1 || k > 1000 {
			fmt.Println("Количество сотрудников от 1 до 1000")

			return
	}

	temperature(k)

	}
}