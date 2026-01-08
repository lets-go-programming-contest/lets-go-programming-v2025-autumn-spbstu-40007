package main

import (
	"container/heap"
	"fmt"
	"os"
)

type thisHeap []int

func (h *thisHeap) Len() int           { return len(*h) }
func (h *thisHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *thisHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *thisHeap) Push(x any) { *h = append(*h, x.(int)) }

func (h *thisHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var count, rank, choice int

	ratingHeap := &thisHeap{}

	_, err := fmt.Fscan(os.Stdin, &count)
	if err != nil {
		fmt.Println("Invalid format")
		return
	}

	if !(count >= 1 && count <= 10000) {
		fmt.Println("The number is not included in the range [1, 10000]")
		return
	}

	for range count {

		_, err = fmt.Fscan(os.Stdin, &rank)
		if err != nil {
			fmt.Println("Invalid format")
			return
		}

		if !(rank >= -10000 && rank <= 10000) {
			fmt.Println("The dish rating is not included in the range [-10000, 10000]")
			return
		}

		heap.Push(ratingHeap, rank)
	}

	_, err = fmt.Fscan(os.Stdin, &choice)
	if err != nil {
		fmt.Println("Invalid format")
		return
	}

	if !(choice >= 1 && choice <= count) {
		fmt.Printf("The number is not included in the range [1, %d]", count)
		return
	}

	for i := 1; i <= (count - choice); i++ {
		heap.Pop(ratingHeap)
	}
	fmt.Println(heap.Pop(ratingHeap))

}
