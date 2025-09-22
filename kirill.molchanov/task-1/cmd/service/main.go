package main

import (
	"fmt"
)

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

func multiply(a, b int) int {
	return a * b
}

func divide(a, b int) int { return a / b }

func main() {
	var (
		a, b      int
		operation string
	)

	_, err := fmt.Scanln(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scanln(&operation)
	if err != nil {
		fmt.Println("invalid input")
		return
	}

	switch operation {
	case "+":
		fmt.Println(add(a, b))
	case "-":
		fmt.Println(subtract(a, b))
	case "*":
		fmt.Println(multiply(a, b))
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(divide(a, b))
	default:
		fmt.Println("Invalid operation")
	}
}
