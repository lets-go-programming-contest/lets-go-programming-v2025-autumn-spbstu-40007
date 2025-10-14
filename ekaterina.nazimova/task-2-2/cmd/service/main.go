package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
    if value, ok := x.(int); ok {
        *h = append(*h, value)
    }
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var dishAmount int

	_, err := fmt.Scan(&dishAmount)
	if err != nil {
		fmt.Println(err)

		return
	}

	dishHeap := &IntHeap{}
	heap.Init(dishHeap)

	for range make([]int, dishAmount) {
		var dish int

		_, err = fmt.Scan(&dish)
		if err != nil {
			fmt.Println(err)

			return
		}

		heap.Push(dishHeap, dish)
	}

	var dishNumber int

	_, err = fmt.Scan(&dishNumber)
	if err != nil {
		fmt.Println(err)

		return
	}

	var kthDish int

	for range make([]int, dishNumber) {
		val, ok := heap.Pop(dishHeap).(int)

		if !ok {
			fmt.Println("Unexpected type")

			return
		}

		kthDish = val
	}

	fmt.Println(kthDish)
}
