package main

import (
	"container/heap"
	"fmt"
)

type IntMaxHeap []int

func (myHeap IntMaxHeap) Len() int {
	return len(myHeap)
}
func (myHeap IntMaxHeap) Less(i, j int) bool {
	return myHeap[i] > myHeap[j]
}
func (myHeap IntMaxHeap) Swap(i, j int) {
	myHeap[i], myHeap[j] = myHeap[j], myHeap[i]
}
func (myHeap *IntMaxHeap) Push(x interface{}) {
	*myHeap = append(*myHeap, x.(int))
}
func (myHeap *IntMaxHeap) Pop() interface{} {
	old := *myHeap
	n := len(old)
	x := old[n-1]
	*myHeap = old[0 : n-1]
	return x
}

func main() {
	var count, value, numberOfDish int

	h := &IntMaxHeap{}

	_, err := fmt.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < count; i++ {
		_, err = fmt.Scan(&value)
		if err != nil {
			fmt.Println(err)
		}

		heap.Push(h, value)
	}
	heap.Init(h)

	_, err = fmt.Scan(&numberOfDish)
	if err != nil {
		fmt.Println(err)
	}

	for i := 1; i < numberOfDish; i++ {
		heap.Pop(h)
	}

	fmt.Println(heap.Pop(h))
}
