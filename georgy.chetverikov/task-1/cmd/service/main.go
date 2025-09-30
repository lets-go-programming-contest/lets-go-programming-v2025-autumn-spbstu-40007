package main

import (
	"errors"
	"fmt"
)

func substract(x int, y int) int {
	return x - y
}

func summarize(x int, y int) int {
	return x + y
}

func multiply(x int, y int) int {
	return x * y
}

func divide(x int, y int) (int, error) {
	if y == 0 {
		return 0, errors.New("Division by zero")
	}

	return x / y, nil
}

func main() {
	var x, y int

	_, err := fmt.Scanln(&x)

	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&y)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var option string

	_, err = fmt.Scanln(&option)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch option {
	case "-":
		fmt.Println(substract(x, y))
	case "+":
		fmt.Println(summarize(x, y))
	case "*":
		fmt.Println(multiply(x, y))
	case "/":
		result, err := divide(x, y)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(result)
	default:
		fmt.Println("Invalid operation")
	}
}
