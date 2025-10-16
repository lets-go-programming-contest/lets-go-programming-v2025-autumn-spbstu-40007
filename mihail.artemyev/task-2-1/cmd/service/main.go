package main

import (
	"fmt"
)

func main() {
	var N int // так в задании написано что уж

	if _, err := fmt.Scan(&N); err != nil {
		fmt.Println("Input error: failed to read N (количество отделов) -", err)

		return
	}

	for i := 0; i < N; i++ {
		var kCount int // а это не хочу просто K

		if _, err := fmt.Scan(&kCount); err != nil {
			fmt.Println("Input error: failed to read K (количество сотрудников) -", err)

			return
		}

		minT := 15
		maxT := 30

		for j := 0; j < kCount; j++ {
			var (
				operator string // оператор сравнения 
				valueT   int
			)

			if _, err := fmt.Scan(&operator, &valueT); err != nil {
				fmt.Println("Input error: failed to read operator and valueT (данные температуры) -", err)

				return
			}

			if operator == ">=" && valueT > minT {
				minT = valueT
			}
			if operator == "<=" && valueT < maxT {
				maxT = valueT
			}

			if minT > maxT {
				fmt.Println(-1)

				continue
			}

			fmt.Println(minT)
		}
	}
}
