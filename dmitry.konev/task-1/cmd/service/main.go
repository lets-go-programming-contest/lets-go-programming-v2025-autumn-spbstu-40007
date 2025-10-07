package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the first operand:")
	if !scanner.Scan() {
		fmt.Println("Invalid first operand")
		return
	}
	firstInput := scanner.Text()
	operand1, err := strconv.Atoi(strings.TrimSpace(firstInput))
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	fmt.Println("Enter second operand:")
	if !scanner.Scan() {
		fmt.Println("Invalid second operand")
		return
	}
	secondInput := scanner.Text()
	operand2, err := strconv.Atoi(strings.TrimSpace(secondInput))
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Println("Enter operation (+, -, *, /):")
	if !scanner.Scan() {
		fmt.Println("Invalid operation")
		return
	}
	op := strings.TrimSpace(scanner.Text())

	switch op {
	case "+":
		fmt.Println(operand1 + operand2)
	case "-":
		fmt.Println(operand1 - operand2)
	case "*":
		fmt.Println(operand1 * operand2)
	case "/":
		if operand2 == 0 {
			fmt.Println("Division by zero")
		} else {
			result := float64(operand1) / float64(operand2)
			fmt.Println(result)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
