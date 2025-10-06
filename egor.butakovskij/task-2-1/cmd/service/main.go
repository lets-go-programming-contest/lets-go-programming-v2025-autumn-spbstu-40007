package main

import "fmt"

func getRecomendedTemperature(K int) {
	var (
		temp, recTemp, highBorder, lowBorder int
		sign                                 string
	)
	recTemp = 15
	highBorder = 30
	lowBorder = 15

	for j := 0; j < K; j++ {
		_, err := fmt.Scanf("%s %d", &sign, &temp)
		if err != nil {
			fmt.Println(err)
		}

		if sign == ">=" {
			lowBorder = max(lowBorder, temp)
		}

		if sign == "<=" {
			highBorder = min(highBorder, temp)
		}

		if lowBorder > highBorder {
			recTemp = -1
		} else {
			recTemp = lowBorder
		}

		fmt.Println(recTemp)
	}
}

func main() {
	var (
		N, K int
	)

	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < N; i++ {
		_, err := fmt.Scan(&K)
		if err != nil {
			fmt.Println(err)
		}
		getRecomendedTemperature(K)
	}
}
