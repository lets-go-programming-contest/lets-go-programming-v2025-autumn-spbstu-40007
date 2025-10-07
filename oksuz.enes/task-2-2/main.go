package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (heap MinHeap) Len() int {
	return len(heap)
}

func (heap MinHeap) Less(i, j int) bool {
	return (heap[i] < heap[j])
}

func (heap MinHeap) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]

}

func (heap *MinHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*heap = append(*heap, value)
}

func (heap *MinHeap) Pop() interface{} {
	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]

	return x
}

func main() {
	var n, k int

	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	prefer := make([]int, n)
	for i := range prefer {
		_, err := fmt.Scan(&prefer[i])
		if err != nil {
			return
		}
	}

	if _, err := fmt.Scan(&k); err != nil {
		return
	}

	if k <= 0 || n == 0 {
		return
	}

	myHeap := &MinHeap{}
	heap.Init(myHeap)

	for _, pref := range prefer {
		if myHeap.Len() < k {
			heap.Push(myHeap, pref)
		} else if pref > (*myHeap)[0] {
			heap.Pop(myHeap)
			heap.Push(myHeap, pref)
		}
	}

	fmt.Println((*myHeap)[0])

}
