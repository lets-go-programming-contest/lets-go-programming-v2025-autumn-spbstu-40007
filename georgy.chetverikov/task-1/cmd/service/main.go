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
		return 0, errors.New("1") 
	}

	return x / y, nil
}

func main() {

	var x, y int

	_, err1 := fmt.Scanln(&x)

	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err2 := fmt.Scanln(&y)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var option string

	_, err3 := fmt.Scanln(&option)
	if err3 != nil {
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
			fmt.Println("Divison by zero")
			return
		} else {
			fmt.Println(result)
		}
	default:
		fmt.Println("Invalid operation")
		return
	}

}
