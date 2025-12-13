package main

import (
	"container/heap"
	"fmt"
)

type intHeap []int

func (height *intHeap) Len() int           { return len(*height) }
func (height *intHeap) Less(i, j int) bool { return (*height)[i] > (*height)[j] }
func (height *intHeap) Swap(i, j int)      { (*height)[i], (*height)[j] = (*height)[j], (*height)[i] }

func (height *intHeap) Push(x interface{}) {
	if value, ok := x.(int); ok {
		*height = append(*height, value)
	}
}

func (height *intHeap) Pop() interface{} {
	old := *height
	n := len(old)
	x := old[n-1]
	*height = old[0 : n-1]

	return x
}

func main() {
	var dishAmount int

	_, err := fmt.Scan(&dishAmount)
	if err != nil {
		fmt.Println(err)

		return
	}

	dishHeap := &intHeap{}
	*dishHeap = make([]int, 0, dishAmount)
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

	var dishNum int

	_, err = fmt.Scan(&dishNum)
	if err != nil {
		fmt.Println(err)

		return
	}

	var kthDish interface{}

	for range make([]int, dishNum) {
		kthDish = heap.Pop(dishHeap)
	}

	fmt.Println(kthDish)
}
