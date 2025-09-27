package main

import (
	"fmt"
)

func calc(x int, y int, z string) {
	switch z {
	case "+":
		fmt.Println(x + y)
	case "-":
		fmt.Println(x - y)
	case "*":
		fmt.Println(x * y)
	case "/":
		if y == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(x / y)
	default:
		fmt.Println("Invalid operation")
	}
}

func main() {
	var (
		a, b int
		c    string
	)

	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&c)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	calc(a, b, c)
}
