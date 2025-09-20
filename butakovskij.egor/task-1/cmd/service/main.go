package main

import (
	"fmt"
	"os"
)

func calc(x int, y int, z string) {
	switch z {
	case "+":
		fmt.Println(x + y)
	case "-":
		fmt.Println(x - y)
	case "*":
		fmt.Println(x * y)
	case "/":
		if y == 0 {
			fmt.Println("Division by zero")
			os.Exit(1)
		} else {
			fmt.Println(x / y)
		}
	}
}

func main() {
	var a int
	var b int
	var c string
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		os.Exit(1)
	}
	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		os.Exit(1)
	}
	_, err = fmt.Scan(&c)
	if err != nil {
		fmt.Println("Invalid operation")
		os.Exit(1)
	}
	if c != "+" && c != "-" && c != "*" && c != "/" {
		fmt.Println("Invalid operation")
		os.Exit(1)
	}

	calc(a, b, c)
}
