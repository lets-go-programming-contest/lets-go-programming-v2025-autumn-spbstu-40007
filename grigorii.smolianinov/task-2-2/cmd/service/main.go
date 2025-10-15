package main

import (
	"container/heap"
	"fmt"
	"log"
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
	var N int //nolint:varnamelen
	if _, err := fmt.Scan(&N); err != nil {
		if err.Error() != "EOF" {
			log.Fatal(err)
		}

		return
	}

	preferences := make([]int, N)
	for i := 0; i < N; i++ {
		if _, err := fmt.Scan(&preferences[i]); err != nil {
			if err.Error() != "EOF" {
				log.Fatal(err)
			}

			return
		}
	}

	var K int //nolint:varnamelen
	if _, err := fmt.Scan(&K); err != nil {
		if err.Error() != "EOF" {
			log.Fatal(err)
		}

		return
	}

	if K < 1 || K > N {
		return
	}

	h := &IntHeap{} //nolint:varnamelen

	for _, pref := range preferences {
		heap.Push(h, pref)

		if h.Len() > K {
			heap.Pop(h)
		}
	}

	if h.Len() > 0 {
		fmt.Println((*h)[0])
	}
}
