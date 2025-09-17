package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	input, _ := reader.ReadString('\n')
	parts := strings.Fields(input)
	if len(parts) > 1 {
		fmt.Println("Enter exactly one numbers")
		return
	}
	if len(parts) == 0 {
		fmt.Println("Empty string")
		return
	}

	a, err1 := strconv.ParseInt(parts[0], 10, 64)
	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	input2, _ := reader.ReadString('\n')
	parts2 := strings.Fields(input2)
	if len(parts2) > 1 {
		fmt.Println("Enter exactly one numbers")
		return
	}
	if len(parts2) == 0 {
		fmt.Println("Empty string")
		return
	}

	b, err2 := strconv.ParseInt(parts2[0], 10, 64)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	opInput, _ := reader.ReadString('\n')
	sign := strings.TrimSpace(opInput)
	if len(sign) > 1 {
		fmt.Println("Enter exactly one operation")
		return
	}
	if len(sign) == 0 {
		fmt.Println("Empty string")
		return
	}

	result, err := calculate(a, b, sign)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

func calculate(n1, n2 int64, s string) (string, error) {
	var result string
	var err error

	switch s {
	case "+":
		result = fmt.Sprintf("%d", n1+n2)
	case "-":
		result = fmt.Sprintf("%d", n1-n2)
	case "*":
		result = fmt.Sprintf("%d", n1*n2)
	case "/":
		if n2 == 0 {
			err = fmt.Errorf("Division by zero")
		} else {
			result = fmt.Sprintf("%d", n1/n2)
		}
	case ":":
		if n2 == 0 {
			err = fmt.Errorf("Divisor by zero")
		} else {
			result = fmt.Sprintf("%d", n1/n2)
		}
	default:
		err = fmt.Errorf("Invalid operation")
	}

	return result, err
}
