package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func die(message string) {
	fmt.Fprintln(os.Stderr, message)

	os.Exit(0)
}

func read(scanner *bufio.Scanner) string {
	scanner.Scan()
	if scanner.Err() != nil {
		die("Quel dommage")
	}

	return scanner.Text()
}

var operations = map[string]func(int64, int64) int64{
	"+": func(x, y int64) int64 { return x + y },
	"-": func(x, y int64) int64 { return x - y },
	"*": func(x, y int64) int64 { return x * y },
	"/": func(x, y int64) int64 { return x / y },
}

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	x, err := strconv.ParseInt(read(stdin), 10, 0)
	if err != nil {
		die("Invalid first operand")
	}

	y, err := strconv.ParseInt(read(stdin), 10, 0)
	if err != nil {
		die("Invalid first operand")
	}

	operation := read(stdin)
	if y == 0 && operation == "/" {
		die("Division by zero")
	}

	if f, ok := operations[operation]; ok {
		fmt.Println(f(x, y))
	} else {
		die("Invalid operation")
	}
}
