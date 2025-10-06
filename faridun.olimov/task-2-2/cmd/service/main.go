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
	if num, ok := x.(int); ok {
		*h = append(*h, num)
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
	var count, kValue int

	_, err := fmt.Scan(&count)
	if err != nil {
		return
	}

	numbers := make([]int, count)

	for i := range count {
		_, err = fmt.Scan(&numbers[i])
		if err != nil {
			return
		}
	}

	_, err = fmt.Scan(&kValue)
	if err != nil {
		return
	}

	intHeap := &IntHeap{}
	heap.Init(intHeap)

	for _, num := range numbers {
		heap.Push(intHeap, num)
	}

	for i := 1; i < kValue; i++ {
		heap.Pop(intHeap)
	}

	fmt.Println((*intHeap)[0])
}
