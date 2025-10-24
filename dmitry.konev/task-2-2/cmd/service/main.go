package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (h *MinHeap) Len() int {
	return len(*h)
}

func (h *MinHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *MinHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MinHeap) Push(x interface{}) {
	xInt, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, xInt)
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[0 : n-1]

	return val
}

func main() {
	var n, kTop int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	arr := make([]int, n)
	for i := range arr {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			return
		}
	}

	if _, err := fmt.Scan(&kTop); err != nil {
		return
	}

	minHeap := &MinHeap{}
	heap.Init(minHeap)

	for _, val := range arr {
		if minHeap.Len() < kTop {
			heap.Push(minHeap, val)
		} else if val > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, val)
		}
	}

	fmt.Println((*minHeap)[0])
}
