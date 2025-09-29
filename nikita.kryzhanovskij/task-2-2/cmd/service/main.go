package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type HeapInterface interface {
	sort.Interface
	Push(x interface{})
	Pop() interface{}
}

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
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

	for i := 0; i < quantity; i++ {
		var dish int

		if _, err := fmt.Scan(&dish); err != nil {
			fmt.Println("Invalid number")

			return
		}

		heap.Push(dishes, dish)
	}

	var k int

	if _, err := fmt.Scan(&k); err != nil || dishes.Len() < k {
		fmt.Println("Invalid number")

		return
	}

	chosen := 0
	for i := 0; i < k; i++ {
		chosen = heap.Pop(dishes).(int)
	}

	fmt.Println(chosen)
}
