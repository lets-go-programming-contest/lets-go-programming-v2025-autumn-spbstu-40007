package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (minHeap MinHeap) Len() int {
	return len(minHeap)
}

func (minHeap MinHeap) Less(i, j int) bool {
	return minHeap[i] < minHeap[j]
}

func (minHeap MinHeap) Swap(i, j int) {
	minHeap[i], minHeap[j] = minHeap[j], minHeap[i]
}

func (minHeap *MinHeap) Push(value interface{}) {
	intValue, ok := value.(int)
	if !ok {
		return
	}
	*minHeap = append(*minHeap, intValue)
}

func (minHeap *MinHeap) Pop() interface{} {
	old := *minHeap
	length := len(old)
	result := old[length-1]
	*minHeap = old[0 : length-1]

	return result
}

func main() {
	var count int

	_, err := fmt.Scan(&count)
	if err != nil || count < 1 || count > 10000 {
		return
	}

	numbers := make([]int, count)

	for idx := range count {
		_, readErr := fmt.Scan(&numbers[idx])
		if readErr != nil {
			return
		}
	}

	var position int

	_, err = fmt.Scan(&position)
	if err != nil || position < 1 || position > count {
		return
	}

	minHeap := &MinHeap{}
	heap.Init(minHeap)

	for idx := range position {
		heap.Push(minHeap, numbers[idx])
	}

	for idx := position; idx < count; idx++ {
		if numbers[idx] > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, numbers[idx])
		}
	}

	fmt.Println((*minHeap)[0])
}
