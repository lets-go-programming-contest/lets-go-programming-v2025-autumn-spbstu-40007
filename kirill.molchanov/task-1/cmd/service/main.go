package main

import (
	"fmt"
	"os"
)

func add(a, b int) int {
	return a + b
}

func subtraction(a, b int) int {
	return a - b
}

func increase(a, b int) int {
	return a * b
}

func division(a, b int) (int, error) {
	if b == 0 {
		fmt.Fprintf(os.Stderr, "division by zero")
		os.Exit(1)
	}
	return a / b, nil
}

func main() {
	var a, b int
	var index string

	_, err := fmt.Scanln(&a)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid first operand")
		os.Exit(1)
	}

	_, err = fmt.Scanln(&b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid second operand")
		os.Exit(1)
	}

	_, err = fmt.Scanf("%s", &index)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid input")
		os.Exit(1)
	}

	if index == "+" {
		fmt.Println(add(a, b))
	} else if index == "-" {
		fmt.Println(subtraction(a, b))
	} else if index == "*" {
		fmt.Println(increase(a, b))
	} else if index == "/" {
		fmt.Println(division(a, b))
	} else {
		fmt.Fprintln(os.Stderr, "Invalid operation")
		os.Exit(1)
	}
}
