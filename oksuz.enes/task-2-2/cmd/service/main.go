package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (myHeap *MinHeap) Len() int {
	return len(*myHeap)
}

func (myHeap *MinHeap) Less(i, j int) bool {
	return (*myHeap)[i] < (*myHeap)[j]
}

func (myHeap *MinHeap) Swap(i, j int) {
	(*myHeap)[i], (*myHeap)[j] = (*myHeap)[j], (*myHeap)[i]
}

func (myHeap *MinHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*myHeap = append(*myHeap, value)
}

func (myHeap *MinHeap) Pop() interface{} {
	old := *myHeap
	n := len(old)
	x := old[n-1]
	*myHeap = old[0 : n-1]

	return x
}

func main() {
	var count, value int

	_, err := fmt.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	prefer := make([]int, count)
	for i := range prefer {
		_, err = fmt.Scan(&prefer[i])
		if err != nil {
			fmt.Println(err)
		}
	}

	_, err = fmt.Scan(&value)
	if err != nil {
		fmt.Println(err)
	}

	minHeap := &MinHeap{}
	heap.Init(minHeap)

	for _, pref := range prefer {
		if minHeap.Len() < value {
			heap.Push(minHeap, pref)
		} else if pref > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, pref)
		}
	}

	fmt.Println((*minHeap)[0])

}
