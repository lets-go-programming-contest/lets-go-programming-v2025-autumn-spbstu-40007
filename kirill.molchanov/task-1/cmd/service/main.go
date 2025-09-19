package main

import (
	"fmt"
	"os"
)

func add(a, b int) int {
	return a + b
}

func subtraction(a, b int) int {
	return a - b
}

func increase(a, b int) int {
	return a * b
}

func division(a, b int) int {
	if b == 0 {
		fmt.Println("Division by zero")
		os.Exit(0)
	}
	return a / b
}

func main() {
	var a, b int
	var index string

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

	_, err = fmt.Scanf("%s", &index)
	if err != nil {
		fmt.Println("invalid input")
		return
	}

	if index == "+" {
		fmt.Println(add(a, b))
	} else if index == "-" {
		fmt.Println(subtraction(a, b))
	} else if index == "*" {
		fmt.Println(increase(a, b))
	} else if index == "/" {
		fmt.Println(division(a, b))
	} else {
		fmt.Println("Invalid operation")
		return
	}
}
