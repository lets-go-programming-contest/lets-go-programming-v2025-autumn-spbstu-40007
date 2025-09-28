package main

// Newline for separation std packages from others seems crazy to linter.
//nolint:gofumpt
import (
	"fmt"
	"strconv"

	"task-1/internal/die"
	"task-1/internal/scanner"
)

var operations = map[string]func(int, int) int{
	"+": func(x, y int) int { return x + y },
	"-": func(x, y int) int { return x - y },
	"*": func(x, y int) int { return x * y },
	"/": func(x, y int) int { return x / y },
}

func main() {
	scanner := scanner.NewScanner()
	x, err := strconv.Atoi(scanner.Read())
	if err != nil {
		die.Die("Invalid first operand")
	}

	y, err := strconv.Atoi(scanner.Read())
	if err != nil {
		die.Die("Invalid second operand")
	}

	operation := scanner.Read()
	if y == 0 && operation == "/" {
		die.Die("Division by zero")
	}

	if f, ok := operations[operation]; ok {
		fmt.Println(f(x, y))
	} else {
		die.Die("Invalid operation")
	}
}
