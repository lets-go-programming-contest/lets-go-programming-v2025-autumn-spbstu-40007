package main

import "fmt"

func main() {
	var (
		num1     float64
		num2     float64
		operator string
		err      error
	)

	_, err = fmt.Scanln(&num1)
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

	if operator != "+" && operator != "-" && operator != "*" && operator != "/" {
		fmt.Println("Invalid operation")
		return
	}

	if num2 == 0 && operator == "/" {
		fmt.Println("Division by zero")
		return
	}

	var result float64
	switch operator {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		result = num1 / num2
	}

	fmt.Println(result)
}
