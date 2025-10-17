package main

import (
	"container/heap"
	"fmt"
)

type stdHeap []int

func (ourHeap stdHeap) Len() int {
	return len(ourHeap)
}

func (ourHeap stdHeap) Less(i, j int) bool {
	return ourHeap[i] < ourHeap[j]
}

func (ourHeap stdHeap) Swap(i, j int) {
	ourHeap[i], ourHeap[j] = ourHeap[j], ourHeap[i]
}

func (ourHeap *stdHeap) Push(x any) {
	*ourHeap = append(*ourHeap, x.(int))
}

func (ourHeap *stdHeap) Pop() any {
	old := *ourHeap
	n := len(old)
	x := old[n-1]
	*ourHeap = old[:n-1]

	return x
}

func main() {
	var numOfNumbers, indexK int

	if _, err := fmt.Scan(&numOfNumbers); err != nil {
		fmt.Println(-1)

		return
	}

	numberOfMeals := make([]int, numOfNumbers)

	for i := 0; i < numOfNumbers; i++ {
		fmt.Scan(&numberOfMeals[i])
	}

	if _, err := fmt.Scan(&indexK); err != nil {
		fmt.Println(-1)

		return
	}

	ourHeap := &stdHeap{}

	heap.Init(ourHeap)

	for i := 0; i < indexK; i++ {
		heap.Push(ourHeap, numberOfMeals[i])
	}

	for i := indexK; i < numOfNumbers; i++ {
		if numberOfMeals[i] > (*ourHeap)[0] {
			heap.Pop(ourHeap)
			heap.Push(ourHeap, numberOfMeals[i])
		}
	}
	fmt.Println((*ourHeap)[0])
}
