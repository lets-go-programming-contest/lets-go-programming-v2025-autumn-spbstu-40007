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

	myHeap := &IntMaxHeap{}

	_, err := fmt.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	for range count {
		_, err = fmt.Scan(&value)
		if err != nil {
			fmt.Println(err)
		}

		heap.Push(myHeap, value)
	}
	
	heap.Init(myHeap)

	_, err = fmt.Scan(&numberOfDish)
	if err != nil {
		fmt.Println(err)
	}

	for range numberOfDish {
		heap.Pop(myHeap)
	}

	fmt.Println(heap.Pop(myHeap))
}
