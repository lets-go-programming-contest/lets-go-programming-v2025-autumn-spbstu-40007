package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

//nolint:recvcheck
type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x any) {
	//nolint:forcetypeassert
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var numberDishes int
	_, err := fmt.Fscan(reader, &numberDishes)

	if err != nil || numberDishes < 1 || numberDishes > 10000 {
		fmt.Println("Invalid dishes range")

		return
	}

	sequenceNumbers := make([]int, numberDishes)

	for i := range sequenceNumbers {
		_, err := fmt.Fscan(reader, &sequenceNumbers[i])

		if err != nil || sequenceNumbers[i] < -10000 || sequenceNumbers[i] > 10000 {
			fmt.Println("Invalid ai range")

			return
		}
	}

	var kDishes int
	_, err2 := fmt.Fscan(reader, &kDishes)

	if err2 != nil || kDishes < 1 || kDishes > numberDishes {
		fmt.Println("Invalid k range")

		return
	}

	heapD := &IntHeap{}
	heap.Init(heapD)

	for i := range sequenceNumbers {
		heap.Push(heapD, sequenceNumbers[i])

		if heapD.Len() > kDishes {
			heap.Pop(heapD)
		}
	}

	fmt.Println((*heapD)[0])
}
