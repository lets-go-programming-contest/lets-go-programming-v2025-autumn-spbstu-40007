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
			fmt.Println("Error:", err)
		}

		if temp < lowBorder && sign == "<=" || temp > highBorder && sign == ">=" {
			fmt.Println("-1")
			return
		}
		if sign == ">=" && temp <= highBorder && temp >= lowBorder {
			recTemp = temp
			lowBorder = temp
		}

		if sign == "<=" && temp <= highBorder && temp >= lowBorder {
			highBorder = temp
		}

		if sign == "<=" && recTemp > temp {
			fmt.Println("-1")
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
		fmt.Println("Error:", err)
	}

	for i := 0; i < N; i++ {
		_, err := fmt.Scan(&K)
		if err != nil {
			fmt.Println("Error:", err)
		}
		getRecomendedTemperature(K)
	}
}
