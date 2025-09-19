package main

import "fmt"

func main() {
	var firstOperand, secondOperand int
	var operation string

	if n, err := fmt.Scan(&firstOperand, &secondOperand, &operation); err != nil {
		switch n {
		case 0:
			fmt.Println("Invalid first operand")
		case 1:
			fmt.Println("Invalid second operand")
		case 2:
			fmt.Println("Invalid operation")
		}
		return
	}

	switch operation {
	case "+":
		fmt.Println(firstOperand + secondOperand)
	case "-":
		fmt.Println(firstOperand - secondOperand)
	case "*":
		fmt.Println(firstOperand * secondOperand)
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
			return
		}

		fmt.Println(firstOperand / secondOperand)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
