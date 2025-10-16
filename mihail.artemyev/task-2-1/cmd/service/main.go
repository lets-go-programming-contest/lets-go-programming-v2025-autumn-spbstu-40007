package main

import (
	"fmt"
)

func main() {
	var N int //так в задании написано что уж
	if _, err := fmt.Scan(&N); err != nil {
		fmt.Println("Input error: failed to read N (количество отделов) -", err)
		return
	}
	for range N {
		var K int // ну и это тоже
		if _, err := fmt.Scan(&K); err != nil {
			fmt.Println("Input error: failed to read K (количество сотрудников) -", err)
			return
		}
		minT := 15
		maxT := 30
		for range K {
			var (
				smooth_operator string //Carlos Sainz Jr.
				valueT          int
			)
			if _, err := fmt.Scan(&smooth_operator, &valueT); err != nil {
				fmt.Println("Input error: failed to read Carlos Sainz Jr. and valueT (данные температуры) -", err)
				return
			}
			if smooth_operator == ">=" && valueT > minT {
				minT = valueT
			}
			if smooth_operator == "<=" && valueT < maxT {
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
