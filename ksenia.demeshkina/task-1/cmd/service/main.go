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
			fmt.Println("Division by zero.")
			return //чтобы не было panic
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
		var clean string
		fmt.Fscanln(os.Stdin, &clean) //очищаем буфер после каждой ошибки
		return
	}

	_, err = fmt.Fscan(os.Stdin, &number2)
	if err != nil {
		fmt.Println("Invalid second operand")
		var clean string
		fmt.Fscanln(os.Stdin, &clean)
		return
	}

	_, err = fmt.Fscan(os.Stdin, &symbol)
	if err != nil {
		fmt.Println("Invalid operation") //техническая ошибка
		var clean string
		fmt.Fscanln(os.Stdin, &clean)
		return
	}

	if symbol != "+" && symbol != "-" && symbol != "/" && symbol != "*" { //логическая ошибка
		fmt.Println("Invalid operation")
		return
	}

	calculate(number1, number2, symbol)
}
