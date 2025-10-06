package main

import (
	"container/heap"
	"fmt"
)

type DishHeap []int

func (h DishHeap) Len() int           { return len(h) }
func (h DishHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h DishHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *DishHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *DishHeap) Pop() any {
	heap := *h
	n := len(heap)
	if n == 0 {
		return nil
	}
	last := heap[n-1]
	*h = heap[:n-1]
	return last
}

func main() {
	var n, k int
	if _, err := fmt.Scan(&n); err != nil || n < 1 || n > 10000 {
		fmt.Println("Invalid input")
		return
	}

	scores := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Scan(&scores[i]); err != nil || scores[i] < -10000 || scores[i] > 10000 {
			fmt.Println("Invalid input")
			return
		}
	}

	if _, err := fmt.Scan(&k); err != nil || k < 1 || k > n {
		fmt.Println("Invalid input")
		return
	}

	h := &DishHeap{}
	heap.Init(h)
	for _, score := range scores {
		heap.Push(h, score)
	}

	for i := 0; i < k-1; i++ {
		heap.Pop(h)
	}

	if h.Len() > 0 {
		fmt.Println((*h)[0])
	} else {
		fmt.Println("Invalid input")
	}
}
