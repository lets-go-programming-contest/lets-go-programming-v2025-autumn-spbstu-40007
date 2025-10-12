package main

import (
	"fmt"
)

func main() {
	var (
		s1, s2    int
		operation string
	)

	if _, err := fmt.Scanln(&s1); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scanln(&s2); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, err := fmt.Scanln(&operation); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(s1 + s2)
	case "-":
		fmt.Println(s1 - s2)
	case "*":
		fmt.Println(s1 * s2)
	case "/":
		if s2 == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(s1 / s2)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
