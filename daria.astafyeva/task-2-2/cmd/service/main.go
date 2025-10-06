package main

import (
	"container/heap"
	"fmt"
)

type DishHeap []int

func (h DishHeap) Len() int           { return len(h) }
func (h DishHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h DishHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *DishHeap) Push(x any) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *DishHeap) Pop() any {
	heapArr := *h
	n := len(heapArr)

	if n == 0 {
		return nil
	}

	last := heapArr[n-1]
	*h = heapArr[:n-1]

	return last
}

func main() {
	var numDishes, preferenceRank int
	if _, err := fmt.Scan(&numDishes); err != nil || numDishes < 1 || numDishes > 10000 {
		fmt.Println("Invalid input")

		return
	}

	scores := make([]int, numDishes)
	for i := range numDishes {
		if _, err := fmt.Scan(&scores[i]); err != nil || scores[i] < -10000 || scores[i] > 10000 {
			fmt.Println("Invalid input")

			return
		}
	}

	if _, err := fmt.Scan(&preferenceRank); err != nil || preferenceRank < 1 || preferenceRank > numDishes {
		fmt.Println("Invalid input")
		return
	}

	dishPreferenceHeap := &DishHeap{}
	heap.Init(dishPreferenceHeap)

	for _, score := range scores {
		heap.Push(dishPreferenceHeap, score)
	}

	for range preferenceRank - 1 {
		heap.Pop(dishPreferenceHeap)
	}

	if dishPreferenceHeap.Len() > 0 {
		fmt.Println((*dishPreferenceHeap)[0])
	} else {
		fmt.Println("Invalid input")
	}
}
