package main

import (
	"container/heap"
	"fmt"
	"log"
)

type IntMaxHeap []int

func (h *IntMaxHeap) Len() int {
	return len(*h)
}

func (h *IntMaxHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntMaxHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntMaxHeap) Push(x interface{}) {
	if value, ok := x.(int); ok {
		*h = append(*h, value)
	}
}

func (h *IntMaxHeap) Pop() interface{} {
	if h.Len() == 0 {
		return nil
	}
	n := h.Len() - 1
	x := (*h)[n]
	*h = (*h)[:n]

	return x
}

func main() {
	log.SetFlags(0)

	var count, value, numberOfDish int

	myHeap := &IntMaxHeap{}

	_, err := fmt.Scan(&count)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for range count {
		_, err = fmt.Scan(&value)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		heap.Push(myHeap, value)
	}

	heap.Init(myHeap)

	_, err = fmt.Scan(&numberOfDish)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for i := 1; i < numberOfDish; i++ {
		if myHeap.Len() == 0 {
			log.Fatal("heap is empty")
		}

		heap.Pop(myHeap)
	}

	if myHeap.Len() == 0 {
		log.Fatal("heap is empty")
	}

	fmt.Println(heap.Pop(myHeap))
}
