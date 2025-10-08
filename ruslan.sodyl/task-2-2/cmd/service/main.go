package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func process(dishes uint16) {
	var (
		numb int
		key  uint16
	)

	kucha := &IntHeap{}
	heap.Init(kucha)

	for range dishes {
		_, err := fmt.Scan(&numb)
		if err != nil {
			fmt.Println("Invalid number")

			continue
		}

		heap.Push(kucha, numb)
	}

	_, err := fmt.Scan(&key)
	if err != nil || key < 1 || key > 10000 {
		fmt.Println("Invalid k-number")

		return
	}

	sort.Slice(*kucha, func(i, j int) bool { return (*kucha)[i] < (*kucha)[j] })

	key = dishes - key
	fmt.Println((*kucha)[key])
}

func main() {
	var dishes uint16

	_, err := fmt.Scan(&dishes)
	if err != nil || dishes < 1 || dishes > 10000 {
		fmt.Println("Invalid number of dishes")

		return
	}

	process(dishes)
}


