package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[0 : n-1]
	return val
}

func main() {
	var n, k int
	fmt.Scan(&n)

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i])
	}

	fmt.Scan(&k)

	h := &MinHeap{}
	heap.Init(h)

	for _, val := range a {
		if h.Len() < k {
			heap.Push(h, val)
		} else if val > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, val)
		}
	}

	fmt.Println((*h)[0])
}
