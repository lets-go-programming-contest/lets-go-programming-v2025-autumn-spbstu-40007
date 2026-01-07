package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (m MinHeap) Len() int {
	return len(m)
}

func (m MinHeap) Less(i, j int) bool {
	return m[i] < m[j]
}

func (m MinHeap) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m *MinHeap) Push(val any) {
	*m = append(*m, val.(int))
}

func (m *MinHeap) Pop() any {
	prev := *m
	size := len(prev)
	result := prev[size-1]
	*m = prev[0 : size-1]
	return result
}

func main() {
	var n int
	_, err := fmt.Scanln(&n)
	if err != nil || n < 1 || n > 10000 {
		fmt.Println("ERROR: invalid number of dishes.")
		return
	}

	ratings := make([]int, n)
	for idx := 0; idx < n; idx++ {
		_, err = fmt.Scan(&ratings[idx])
		if err != nil || ratings[idx] < -10000 || ratings[idx] > 10000 {
			fmt.Println("ERROR: invalid rating of dishes.")
			return
		}
	}

	var k int
	_, err = fmt.Scanln(&k)
	if err != nil || k < 1 || k > n {
		fmt.Println("ERROR: invalid number of desired dish")
		return // ← ЭТО ДОБАВЛЕНО!
	}

	minHeap := &MinHeap{}
	heap.Init(minHeap)

	for idx := 0; idx < k; idx++ {
		heap.Push(minHeap, ratings[idx])
	}

	for idx := k; idx < n; idx++ {
		if ratings[idx] > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, ratings[idx])
		}
	}

	fmt.Println((*minHeap)[0])
}
