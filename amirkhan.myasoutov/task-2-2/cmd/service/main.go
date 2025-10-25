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

func readDishAmount() (int, error) {
	var dishAmount int
	if _, err := fmt.Scanln(&dishAmount); err != nil {
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
		var rating int
		if _, err := fmt.Scan(&rating); err != nil {
			return nil, fmt.Errorf("Invalid dish rating")
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
	if _, err := fmt.Scan(&preferenceNumber); err != nil {
		return 0, fmt.Errorf("Invalid number of preference")
	}

	if preferenceNumber < 0 || preferenceNumber > dishAmount {
		return 0, fmt.Errorf("Invalid number of preference")
	}

	return preferenceNumber, nil
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

	dishRating := &DishRating{}
	heap.Init(dishRating)

	for _, value := range dishQueue {
		if dishRating.Len() < preferenceNumber {
			heap.Push(dishRating, value)
			continue
		}

		if dishRating.Len() > 0 && (*dishRating)[0] < value {
			heap.Pop(dishRating)
			heap.Push(dishRating, value)
		}
	}

	if dishRating.Len() == 0 {
		fmt.Println(0)
		return
	}

	result := (*dishRating)[0]
	fmt.Println(result)
}
