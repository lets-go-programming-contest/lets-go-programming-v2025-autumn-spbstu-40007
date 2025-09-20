
package main

import (
	"fmt"
	"errors"
	"log"
)

func substract(x int, y int) int {
	return x - y;	
}

func summarize(x int, y int) int {
	return x + y;
}

func multiply(x int, y int) int {
	return x * y;
}

func divide(x int, y int) (int, error) {
	if y == 0 {
		return 0, errors.New("division by zero");
	}

	return x/y, nil;
}

func main() {
	
	var x,y int;

	_, err1 := fmt.Scanln(&x);

	if err1 != nil {
		log.Fatal(errors.New("invalid fisrt operand"));
	}

	_, err2 := fmt.Scanln(&y);
	if err2 != nil {
		log.Fatal(errors.New("invalid second operand"));
	}

	var option string;

	_, err3 := fmt.Scanln(&option);
	if err3 != nil {
		log.Fatal(errors.New("invalid operation"));
	}

	switch option {
	case "-":
		fmt.Println(substract(x,y));
	case "+":
		fmt.Println(summarize(x,y));
	case "*":
		fmt.Println(multiply(x,y));
	case "/":
		result, err := divide(x,y);
		if  err != nil {
			fmt.Println("Error: ", err);
		} else {
			fmt.Println(result);
		}
	default:
		log.Fatal(errors.New("invalid operation"));	
	}
}


