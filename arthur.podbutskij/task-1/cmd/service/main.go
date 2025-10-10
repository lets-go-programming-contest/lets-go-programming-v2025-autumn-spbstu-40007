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
	input = strings.TrimSpace(input)
	parts := strings.Fields(input)

	first_operand, err_1 := strconv.Atoi(parts[0])
	if (err_1 != nil){
		fmt.Println("Invalid first operand")
		return
	}

	second_operand, err_2 := strconv.Atoi(parts[1])
	if (err_2 != nil){
		fmt.Println("Invalid second operand")
		return
	}

	operation := parts[2]

	var flag bool = true
	var result float64
	switch {
	case operation == "+":
		result = float64(first_operand) + float64(second_operand)
	case operation == "-":
		result = float64(first_operand) - float64(second_operand)
	case operation == "*":
		result = float64(first_operand) * float64(second_operand)
	case operation == "/":
		if (second_operand != 0){
			result = float64(first_operand) / float64(second_operand)
		} else {
			fmt.Println("Division by zero")
			flag = false
		}
	default:
		fmt.Print("Invalid operation")
		flag = false
	}

	if (flag){
		fmt.Println(result)
	} 
}