package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	if value, ok := x.(int); ok {
		*h = append(*h, value)
	}
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func main() {
	var quantity int

	if _, err := fmt.Scan(&quantity); err != nil {
		fmt.Println("Invalid number")

		return
	}

	dishes := &IntHeap{}

	for range quantity {
		var dish int

		if _, err := fmt.Scan(&dish); err != nil {
			fmt.Println("Invalid number")

			return
		}

		heap.Push(dishes, dish)
	}

	var preferred int

	if _, err := fmt.Scan(&preferred); err != nil || dishes.Len() < preferred {
		fmt.Println("Invalid number")

		return
	}

	var chosen int

	for range preferred {
		if value, ok := heap.Pop(dishes).(int); ok {
			chosen = value
		}
	}

	fmt.Println(chosen)
}
