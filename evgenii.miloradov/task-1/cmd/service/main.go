package main

import (
	"fmt"
)

func main() {
	var a int
	var b int
	var operator string

	_, err1 := fmt.Scan(&a)
	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err2 := fmt.Scan(&b)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err3 := fmt.Scan(&operator)
	if err3 != nil {
		fmt.Println("Invalid operator")
		return
	}

	switch operator {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(a / b)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
