package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var dishesNumber int
	_, err := fmt.Scanln(&dishesNumber)
	if err != nil || dishesNumber < 1 || dishesNumber > 10000 {
		fmt.Println("ERROR: invalid number of dishes.")
		return
	}

	dishesRating := make([]int, dishesNumber)
	for i := 0; i < dishesNumber; i++ {
		_, err = fmt.Scan(&dishesRating[i])
		if err != nil || dishesRating[i] < -10000 || dishesRating[i] > 10000 {
			fmt.Println("ERROR: invalid rating of dishes.")
			return
		}
	}

	var desiredDish int
	_, err = fmt.Scanln(&desiredDish)
	if err != nil || desiredDish < 1 || desiredDish > dishesNumber {
		fmt.Println("ERROR: invalid number of desired dish")
	}

	h := &IntHeap{}
	heap.Init(h)

	for i := 0; i < desiredDish; i++ {
		heap.Push(h, dishesRating[i])
	}

	for i := desiredDish; i < dishesNumber; i++ {
		if dishesRating[i] > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, dishesRating[i])
		}
	}

	fmt.Println((*h)[0])
}
