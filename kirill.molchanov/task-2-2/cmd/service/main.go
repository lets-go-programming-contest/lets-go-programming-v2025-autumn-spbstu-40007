package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap struct {
	elements []int
}

func (h *MaxHeap) Len() int {
	return len(h.elements)
}

func (h *MaxHeap) Less(i, j int) bool {
	return h.elements[i] < h.elements[j]
}

func (h *MaxHeap) Swap(i, j int) {
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
}

func (h *MaxHeap) Push(x interface{}) {
	//nolint:forcetypeassert
	h.elements = append(h.elements, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
	old := h.elements
	n := len(old)
	x := old[n-1]
	h.elements = old[0 : n-1]

	return x
}

func (h *MaxHeap) Peek() int {
	return h.elements[0]
}

func main() {
	var countOfMeals, targetIndex int

	_, err := fmt.Scan(&countOfMeals)
	if err != nil {
		fmt.Println(-1)

		return
	}

	numOfMeals := make([]int, countOfMeals)

	//nolint:intrange
	for i := 0; i < countOfMeals; i++ {
		_, err := fmt.Scan(&numOfMeals[i])
		if err != nil {
			fmt.Println(-1)

			return
		}
	}

	_, err = fmt.Scan(&targetIndex)
	if err != nil {
		fmt.Println(-1)

		return
	}

	numHeap := &MaxHeap{
		elements: []int{},
	}

	heap.Init(numHeap)

	//nolint:intrange
	for i := 0; i < targetIndex; i++ {
		heap.Push(numHeap, numOfMeals[i])
	}

	for i := targetIndex; i < countOfMeals; i++ {
		if numOfMeals[i] > numHeap.Peek() {
			heap.Pop(numHeap)
			heap.Push(numHeap, numOfMeals[i])
		}
	}

	fmt.Println(numHeap.Peek())
}
