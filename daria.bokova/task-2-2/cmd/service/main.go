package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap struct {
	elements []int
}

func (h *MaxHeap) Len() int { return len(h.elements) }

func (h *MaxHeap) Less(i, j int) bool {
	return h.elements[i] > h.elements[j]
}

func (h *MaxHeap) Swap(i, j int) {
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
}

func (h *MaxHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		h.elements = append(h.elements, val)
	}
}

func (h *MaxHeap) Pop() interface{} {
	old := h.elements
	n := len(old)
	item := old[n-1]
	h.elements = old[0 : n-1]

	return item
}

func (h *MaxHeap) Peek() int {
	return h.elements[0]
}

func main() {
	var totalNumbers, targetPosition int

	if _, scanErr := fmt.Scan(&totalNumbers); scanErr != nil {
		return
	}

	numberSequence := make([]int, totalNumbers)
	for idx := range totalNumbers {
		if _, scanErr := fmt.Scan(&numberSequence[idx]); scanErr != nil {
			return
		}
	}

	if _, scanErr := fmt.Scan(&targetPosition); scanErr != nil {
		return
	}

	numHeap := &MaxHeap{elements: []int{}}
	heap.Init(numHeap)

	for _, value := range numberSequence {
		heap.Push(numHeap, value)
	}

	for extractCount := 1; extractCount < targetPosition; extractCount++ {
		heap.Pop(numHeap)
	}

	fmt.Println(numHeap.Peek())
}
