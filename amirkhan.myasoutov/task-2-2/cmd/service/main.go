package main

import "fmt"

func main() {
	var (
		N    int
		mass []int
	)

	_, err := fmt.Scanln(&N)
	if err != nil || N < 0 || N > 10000 {
		fmt.Println("Invalid number of dishes")
		return
	}

	for i := 0; i < N; i++ {
		var ai int
		_, err := fmt.Scanln(&ai)
		if err != nil || N < -10000 || N > 10000 {
			fmt.Println("Invalid")
			return
		}
		mass = append(mass, ai)
	}

	for i := 0; i < len(mass); i++ {
		fmt.Println(mass[i])
	}
}
