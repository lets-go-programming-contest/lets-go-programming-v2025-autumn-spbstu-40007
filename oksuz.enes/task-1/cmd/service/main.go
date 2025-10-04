package main

import (
	"fmt"
)

func main() {
	var (
		first, second, result int
		operator              string
	)

	_, err := fmt.Scanln(&first)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&second)
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
		result = first + second
		fmt.Println(result)
	case "-":
		result = first - second
		fmt.Println(result)
	case "*":
		result = first * second
		fmt.Println(result)
	case "/":
		if second == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = first / second
		fmt.Println(result)

	default:
		fmt.Println("Invalid operation")
	}
}
