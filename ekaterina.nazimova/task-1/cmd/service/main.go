package main

import (
	"fmt"
)

func main() {
<<<<<<< HEAD
	var (
		num1, num2 int
		op         string
	)
=======
	var num1, num2 int
	var op string
>>>>>>> 9e86d2a ([TASK-1] add main code)

	if _, err := fmt.Scanln(&num1); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scanln(&num2); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
	if _, err := fmt.Scanln(&op); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch op {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
=======
	fmt.Scanln(&op)
=======
	if _, err := fmt.Scanln(&op); err != nil {
		fmt.Println("Invalid option")
		return
	}
>>>>>>> 00fdeae ([TASK-1] fix error)

	if op == "+" {
		fmt.Println(num1 + num2)
	} else if op == "-" {
		fmt.Println(num1 - num2)
	} else if op == "*" {
		fmt.Println(num1 * num2)
	} else if op == "/" {
>>>>>>> 9e86d2a ([TASK-1] add main code)
		if num2 == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(num1 / num2)
		}
<<<<<<< HEAD
	default:
=======
	} else {
>>>>>>> 9e86d2a ([TASK-1] add main code)
		fmt.Println("Invalid operation")
	}
}
