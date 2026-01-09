package main

import (
	"fmt"
	"os"
)

func main() {
	var num1, num2 int
	var oper string

	_, err := fmt.Fscanln(os.Stdin, &num1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Fscanln(os.Stdin, &num2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Fscanln(os.Stdin, &oper)

	if !((oper == "+") || (oper == "-") || (oper == "*") || (oper == "/")) {
		fmt.Println("Invalid operation")
		return
	}

	if (oper == "/") && (num2 == 0) {
		fmt.Println("Division by zero")
		return
	}

	switch oper {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	default:
		fmt.Println(num1 / num2)
	}

}
