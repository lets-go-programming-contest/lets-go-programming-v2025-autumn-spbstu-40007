package main

import (
	"fmt"
)

func main() {
	var (
		num1, num2 int
		op         string
	)

	if _, err := fmt.Scanln(&num1); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scanln(&num2); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, err := fmt.Scanln(&op); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch op {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
		if num2 == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(num1 / num2)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
