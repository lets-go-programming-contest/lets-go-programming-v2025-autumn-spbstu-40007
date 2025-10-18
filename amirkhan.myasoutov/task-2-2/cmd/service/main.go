package main

import (
	"container/heap"
	"fmt"
	"io"
)

type DishRating []int

func (h *DishRating) Len() int { return len(*h) }
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

func readDishAmount() (int, error) {
	var dishAmount int
	if _, err := fmt.Scanln(&dishAmount); err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("Invalid number of dishes")
		}
		return 0, fmt.Errorf("Invalid number of dishes")
	}
	if dishAmount < 1 || dishAmount > 10000 {
		return 0, fmt.Errorf("Invalid number of dishes")
	}
	return dishAmount, nil
}

func readDishRatings(dishAmount int) ([]int, error) {
	dishQueue := make([]int, dishAmount)
	for index := range dishQueue {
		_, err = fmt.Scan(&dishQueue[index])
		if err != nil {
			fmt.Println("Invalid dish rating")

			return
		}
		if rating < -10000 || rating > 10000 {
			return nil, fmt.Errorf("Invalid dish rating")
		}
		dishQueue[index] = rating
	}
	return dishQueue, nil
}

func readPreferenceNumber(dishAmount int) (int, error) {
	var preferenceNumber int
	_, err = fmt.Scan(&preferenceNumber)
	
	if err != nil {
		fmt.Println("Invalid number of preference")

		return
	}
	if preferenceNumber < 0 || preferenceNumber > dishAmount {
		return 0, fmt.Errorf("Invalid number of preference")
	}
	return preferenceNumber, nil
}

func findNBestRatings(dishQueue []int, preferenceNumber int) int {
	if preferenceNumber == 0 {
		return 0
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
	if dishRating.Len() == 0 {
		return 0
	}
	return (*dishRating)[0]
}

func main() {
	dishAmount, err := readDishAmount()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dishQueue, err := readDishRatings(dishAmount)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	preferenceNumber, err := readPreferenceNumber(dishAmount)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result := findNBestRatings(dishQueue, preferenceNumber)
	fmt.Println(result)
}
