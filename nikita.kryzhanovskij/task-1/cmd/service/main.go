package main

import "fmt"

func main() {
	var first_operand, second_operand int
	var operation string

	if n, err := fmt.Scan(&first_operand, &second_operand, &operation); err != nil {
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
		fmt.Println(first_operand + second_operand)
	case "-":
		fmt.Println(first_operand - second_operand)
	case "*":
		fmt.Println(first_operand * second_operand)
	case "/":
		if second_operand == 0 {
			fmt.Println("Division by zero")
			return
		}

		fmt.Println(first_operand / second_operand)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
