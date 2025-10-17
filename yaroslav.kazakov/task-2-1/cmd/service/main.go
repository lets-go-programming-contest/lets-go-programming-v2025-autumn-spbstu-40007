package main

import (
	"fmt"
)

func main(){
	var N,K int
	_, err := fmt.Scanln(&N)
	if err != nil {
		fmt.Println("Invalid N")
		return
	}
	for i := 0; i < N; i++{
		_, err = fmt.Scanln(&K)
		if err != nil {
			fmt.Println("Invalid K")
			return
		}
		var max = 30
		var min = 15
		for i := 0; i < K; i++{
			var op string
			var T int
			_, err := fmt.Scan(&op, &T)
			if err != nil {
				fmt.Println("invalid Temp input")
				return
			}
			
			switch(op){
				case ">=":
					if T > min{
						min = T
					}
				case "<=":
					if T < max{
						max = T
					}
				default:
					fmt.Println("Invalid: >=/<= expected")
					return
			}
			if min > max {
				fmt.Println(-1)
			}else{
				fmt.Println(min)
			}
		}
	}
	
	
}
