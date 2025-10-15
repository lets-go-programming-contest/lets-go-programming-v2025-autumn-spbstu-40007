package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

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
	var N int
	if _, err := fmt.Scan(&N); err != nil && err.Error() != "EOF" {
		log.Fatal(err)
	}

	preferences := make([]int, N)
	for i := 0; i < N; i++ {
		if _, err := fmt.Scan(&preferences[i]); err != nil && err.Error() != "EOF" {
			log.Fatal(err)
		}
	}

	var K int
	if _, err := fmt.Scan(&K); err != nil && err.Error() != "EOF" {
		log.Fatal(err)
	}

	h := &IntHeap{}

	for _, pref := range preferences {
		heap.Push(h, pref)

		if h.Len() > K {
			heap.Pop(h)
		}
	}

	if h.Len() > 0 {
		fmt.Println((*h)[0])
	} else {
		os.Exit(1)
	}
}
