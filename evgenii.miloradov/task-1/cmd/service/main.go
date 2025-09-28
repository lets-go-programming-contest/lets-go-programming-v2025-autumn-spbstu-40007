package main

import (
	"fmt"
)

func main() {
	var fOperand, sOperand int
	var operator string

	_, err := fmt.Scan(&fOperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&sOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		fmt.Println("Invalid operator")
		return
	}

	switch operator {
	case "+":
		fmt.Println(fOperand + sOperand)
	case "-":
		fmt.Println(fOperand - sOperand)
	case "*":
		fmt.Println(fOperand * sOperand)
	case "/":
		if sOperand == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(fOperand / sOperand)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
