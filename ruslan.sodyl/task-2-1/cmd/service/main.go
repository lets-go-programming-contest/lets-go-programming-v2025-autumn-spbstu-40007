package main
import ("fmt")

func main(){
	var(
		N,K uint16
		temp, minTemp, maxTemp uint8
		sign string
	)
	const(
		minConstTemp = 15
		maxConstTemp = 30
		moreSign = ">="
		lessSign = "<="
	)
  
	_, err := fmt.Scan(&N)
	if err != nil || N > 1000 || N < 1{
		fmt.Println("Invalid number of departmens")
		return
	}

	for range N {
		minTemp = 15
		maxTemp = 30
		_, err = fmt.Scan(&K)
		if err != nil || K > 1000 || K < 1{
			fmt.Println("Invalid number of emploees")
			return
		}
		for range K {
			_, err = fmt.Scan(&sign,&temp)
			if err != nil || sign != lessSign && sign != moreSign || 
			temp > maxConstTemp || temp < minConstTemp {
				fmt.Println("Invalid temperature")
				return
			}
			switch sign {
			case ">=":
				if minTemp < temp {minTemp = temp}
			case "<=":
				if maxTemp > temp {maxTemp = temp}
			}
			if minTemp <= maxTemp{
				fmt.Println(minTemp)
			}else{
				fmt.Println(-1)
			}
		}
	}
}
