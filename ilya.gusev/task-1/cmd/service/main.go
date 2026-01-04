package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		return
	}

	input1 := strings.TrimSpace(sc.Text())
	a, err := strconv.Atoi(input1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if !sc.Scan() {
		return
	}

	input2 := strings.TrimSpace(sc.Text())
	b, err := strconv.Atoi(input2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if !sc.Scan() {
		return
	}

	op := strings.TrimSpace(sc.Text())
	if op == "+" {
		fmt.Println(a + b)
	} else if op == "-" {
		fmt.Println(a - b)
	} else if op == "*" {
		fmt.Println(a * b)
	} else if op == "/" {
		if b == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(a / b)
		}
	} else {
		fmt.Println("Invalid operation")
	}
}
