package main

import (
	"fmt"
	"sort"
)

func main() {
	var dishAmount int
	_, err := fmt.Scan(&dishAmount)

	if err != nil {
		fmt.Println(err)
		return
	}

	dishes := make([]int, dishAmount)
	for i := 0; i < dishAmount; i++ {
		_, err = fmt.Scan(&dishes[i])

		if err != nil {
			fmt.Println(err)
			return
		}

	}

	var dishNumber int
	_, err = fmt.Scan(&dishNumber)

	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Sort(sort.Reverse(sort.IntSlice(dishes)))

	if dishNumber >= 1 && dishNumber <= dishAmount {
		fmt.Println(dishes[dishNumber-1])
	} else {
		fmt.Println("Incorrect dish number")
	}
}
