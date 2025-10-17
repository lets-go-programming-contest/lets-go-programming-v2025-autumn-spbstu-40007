package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int //nolint:recvcheck

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int)) //nolint:forcetypeassert
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() { //nolint:cyclop
	var numOfDishes int
	_, err := fmt.Scanln(&numOfDishes)

	if err != nil || numOfDishes < 1 || numOfDishes > 10000 {
		fmt.Println("Error: invalid number of dishes")

		return
	}

	ratingOfDishes := make([]int, numOfDishes)
	for i := range numOfDishes {
		_, err = fmt.Scan(&ratingOfDishes[i])
		if err != nil || ratingOfDishes[i] < -10000 || ratingOfDishes[i] > 10000 {
			fmt.Println("Error: invalid rating of dishes")

			return
		}
	}

	var dishPriority int
	_, err = fmt.Scan(&dishPriority)

	if err != nil || dishPriority > numOfDishes || dishPriority < 0 {
		fmt.Println("Error: invalid number of priority of current dish")

		return
	}

	h := &IntHeap{} //nolint:varnamelen
	heap.Init(h)

	for i := range dishPriority {
		heap.Push(h, ratingOfDishes[i])
	}

	for i := dishPriority; i < numOfDishes; i++ {
		if ratingOfDishes[i] > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, ratingOfDishes[i])
		}
	}

	fmt.Println((*h)[0])
}
