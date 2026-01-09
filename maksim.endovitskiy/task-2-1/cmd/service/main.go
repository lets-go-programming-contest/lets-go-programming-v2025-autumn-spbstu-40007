package main

import (
	"fmt"
	"os"
)

func main() {
	var department, employee, temp, flag, max, min int
	var opt string
	fmt.Fscanln(os.Stdin, &department)
	for i := 0; i < department; i++ {
		flag = 0
		max = 30
		min = 15
		fmt.Fscanln(os.Stdin, &employee)
		for j := 0; j < employee; j++ {
			fmt.Fscanln(os.Stdin, &opt, &temp)

			if flag == 1 {
				fmt.Println(-1)
				continue
			}

			if opt == ">=" {
				if temp > min {
					min = temp
				}
			} else {
				if temp < max {
					max = temp
				}
			}

			if max < min {
				flag = 1
				fmt.Println(-1)
			} else {
				fmt.Println(min)
			}

		}
	}
}
