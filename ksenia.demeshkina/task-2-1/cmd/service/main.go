package main

import "fmt"

func temperature(K int) {
	minTemp := 15
	maxTemp := 30

	for i := 0; i < K; i++ {
		var operator string;
		var value int;
		fmt.Scan(&operator, &value)

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

			for j := i + 1; j < K; j++ {
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
	var N, K int
	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Println("Error")
	}

	if N < 1 || N > 1000 {
		fmt.Println("Количество отделов в диапазоне от 1 до 1000")
		return
	}

	for i := 0; i < N; i++ {
		_, err = fmt.Scan(&K)
		if err != nil {
			fmt.Println("Error")
	}

		if K < 1 || K > 1000 {
			fmt.Println("Количество сотрудников от 1 до 1000")
			return
	}

	temperature(K)

	}
}