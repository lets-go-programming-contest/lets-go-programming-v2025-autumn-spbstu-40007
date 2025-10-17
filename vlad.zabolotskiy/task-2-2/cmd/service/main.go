package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

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

	a := make([]int, numberDishes)

	for i := range numberDishes {
		_, err := fmt.Fscan(reader, &a[i])

		if err != nil || a[i] < -10000 || a[i] > 10000 {
			fmt.Println("Invalid ai range")

			return
		}
	}

	var k int
	_, err2 := fmt.Fscan(reader, &k)

	if err2 != nil || k < 1 || k > numberDishes {
		fmt.Println("Invalid k range")

		return
	}

	h := &IntHeap{}
	heap.Init(h)

	for i := range a {
		heap.Push(h, a[i])

		if h.Len() > k {
			heap.Pop(h)
		}
	}

	fmt.Println((*h)[0])
}
