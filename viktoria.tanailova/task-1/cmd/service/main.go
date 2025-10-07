package main

import (
	"fmt"
)

func main() {
	var num1, num2 int
	var operator string

	_, err := fmt.Scanln(&num1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&num2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scanln(&operator)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
		if num2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(num1 / num2)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
