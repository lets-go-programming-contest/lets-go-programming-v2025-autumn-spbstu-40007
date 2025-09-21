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

func showErrors(flag1 bool, flag2 bool) bool {
	if flag1 {
		fmt.Println("Invalid first operand")
		return true
	}
	if flag2 {
		fmt.Println("Invalid second operand")
		return true
	}
	return false
}

func main() {
	var flag1, flag2 bool

	reader := bufio.NewScanner(os.Stdin)

	firstNumber, err := strconv.Atoi(readInput(reader))
	if err != nil {
		flag1 = true
	}

	secondNumber, err := strconv.Atoi(readInput(reader))
	if err != nil {
		flag2 = true
	}

	operand := readInput(reader)

	check := showErrors(flag1, flag2)
	if check {
		return
	}

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
