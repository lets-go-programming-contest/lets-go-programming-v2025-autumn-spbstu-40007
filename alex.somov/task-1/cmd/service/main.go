package main

import (
	"fmt"
)

func add(a int, b int) int {
	return a + b
}

func subtract(a int, b int) int {
	return a - b
}

func multiply(a int, b int) int {
	return a * b
}

func divide(a int, b int) int {
	return a / b
}

func main() {
	var firstOperand, secondOperand int
	var operator string

	_, err := fmt.Scanln(&firstOperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&secondOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Scanln(&operator)

	var result int

	switch operator {
	case "+":
		result = add(firstOperand, secondOperand)

	case "-":
		result = subtract(firstOperand, secondOperand)

	case "*":
		result = multiply(firstOperand, secondOperand)

	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
			return
		}

		result = divide(firstOperand, secondOperand)

	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(result)

}
