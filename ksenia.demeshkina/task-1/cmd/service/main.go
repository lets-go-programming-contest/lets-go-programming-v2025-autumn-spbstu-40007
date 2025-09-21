package main

import (
	"fmt"
	"os"
)

func calculate(number1 int, number2 int, symbol string) {
	switch symbol {
	case "+":
		fmt.Println(number1 + number2)
	case "-":
		fmt.Println(number1 - number2)
	case "*":
		fmt.Println(number1 * number2)
	case "/":
		if number2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(number1 / number2)
	}
}

func main() {
	var number1 int
	var number2 int
	var symbol string

	_, err := fmt.Fscan(os.Stdin, &number1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Fscan(os.Stdin, &number2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Fscan(os.Stdin, &symbol)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	if symbol != "+" && symbol != "-" && symbol != "/" && symbol != "*" {
		fmt.Println("Invalid operation")
		return
	}

	calculate(number1, number2, symbol)
}
