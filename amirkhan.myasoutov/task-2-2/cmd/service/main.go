package main

import (
	"container/heap"
	"fmt"
)

type DishRating []int

func (h *DishRating) Len() int           { return len(*h) }
func (h *DishRating) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *DishRating) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *DishRating) Push(x any) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *DishRating) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var dishAmount int
	_, err := fmt.Scanln(&dishAmount)
	
	if err != nil {
		fmt.Println("Invalid number of dishes")

		return
	}

	if dishAmount < 1 || dishAmount > 10000 {
		fmt.Println("Invalid number of dishes")

		return
	}

	dishQueue := make([]int, dishAmount)
	for index := range dishQueue {
		_, err = fmt.Scan(&dishQueue[index])

		if err != nil {
			fmt.Println("Invalid dish rating")

			return
		}

		if dishQueue[index] < -10000 || dishQueue[index] > 10000 {
			fmt.Println("Invalid dish rating")

			return
		}
	}

	var preferenceNumber int
	_, err = fmt.Scan(&preferenceNumber)
	
	if err != nil {
		fmt.Println("Invalid number of preference")

		return
	}

	if preferenceNumber < 0 || preferenceNumber > dishAmount {
		fmt.Println("Invalid number of preference")

		return
	}

	dishRating := &DishRating{}
	heap.Init(dishRating)

	for _, value := range dishQueue {
		if dishRating.Len() < preferenceNumber {
			heap.Push(dishRating, value)

			continue
		}

		if (*dishRating)[0] < value {
			heap.Pop(dishRating)
			heap.Push(dishRating, value)
		}
	}

	result := (*dishRating)[0]
	fmt.Println(result)
}
