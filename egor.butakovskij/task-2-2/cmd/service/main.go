package main

import (
	"container/heap"
	"fmt"
)

type IntMaxHeap []int

func (h IntMaxHeap) Len() int {
	return len(h)
}
func (h IntMaxHeap) Less(i, j int) bool {
	return h[i] > h[j]
}
func (h IntMaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *IntMaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *IntMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var N, ai, k int

	h := &IntMaxHeap{}

	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < N; i++ {
		_, err = fmt.Scan(&ai)
		if err != nil {
			fmt.Println(err)
		}

		heap.Push(h, ai)
	}
	heap.Init(h)

	_, err = fmt.Scan(&k)
	if err != nil {
		fmt.Println(err)
	}

	for i := 1; i < k; i++ {
		heap.Pop(h)
	}

	fmt.Println(heap.Pop(h))
}
