package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readInput(scanner *bufio.Scanner) string {
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		fmt.Println("Error reading input")
		os.Exit(0)
	}
	return scanner.Text()
}

func main() {
	reader := bufio.NewScanner(os.Stdin)

	firstNumber, err := strconv.Atoi(readInput(reader))
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	secondNumber, err := strconv.Atoi(readInput(reader))
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	operand := readInput(reader)

	switch operand {
	case "+":
		fmt.Println(firstNumber + secondNumber)
	case "-":
		fmt.Println(firstNumber - secondNumber)
	case "*":
		fmt.Println(firstNumber * secondNumber)
	case "/":
		if secondNumber == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(firstNumber / secondNumber)
	default:
		fmt.Println("Invalid operation")
	}
}
