package main

import "fmt"

func process(emploeesNumb uint16) {
	var (
		temp 	uint8
		sign 	string
		minTemp uint8 = 15 
		maxTemp uint8 = 30
		
	)
	const (
		minConstTemp = 15
		maxConstTemp = 30
		moreSign = ">="
		lessSign = "<="
	)

	for range emploeesNumb {
		_, err := fmt.Scan(&sign,&temp)
		if err != nil || sign != lessSign && sign != moreSign || 
		temp > maxConstTemp || temp < minConstTemp {
			fmt.Println("Invalid temperature")
			continue 
		}

		switch sign {
		case ">=":
			if minTemp < temp {minTemp = temp}
		case "<=":
			if maxTemp > temp {maxTemp = temp}
		}

		if minTemp <= maxTemp {
			fmt.Println(minTemp)
		}else{
			fmt.Println(-1)
		}
	}
}
func main() {
	var(
		departNumb,emploeesNumb uint16
	)
	
	_, err := fmt.Scan(&departNumb)
	if err != nil || departNumb > 1000 || departNumb < 1 {
		fmt.Println("Invalid number of departments")
		return
	}

	for range departNumb {

		_, err = fmt.Scan(&emploeesNumb)
		if err != nil || emploeesNumb > 1000 || emploeesNumb < 1 {
			fmt.Println("Invalid number of emploees")
			return
		}

		process(emploeesNumb)
	}
}
