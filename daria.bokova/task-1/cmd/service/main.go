package main

import "fmt"

func main() {
	var (
		operand1 float64
		operand2 float64
		operator string
		err      error
	)

	_, err = fmt.Scanln(&operand1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&operand2)
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

	if operand2 == 0 && operator == "/" {
		fmt.Println("Division by zero")
		return
	}

	var res float64
	switch operator {
	case "+":
		res = operand1 + operand2
	case "-":
		res = operand1 - operand2
	case "*":
		res = operand1 * operand2
	case "/":
		res = operand1 / operand2
	}

	fmt.Println(res)
}
