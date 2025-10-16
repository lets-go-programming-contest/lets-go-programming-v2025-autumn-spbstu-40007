package main

import (
	"container/heap"
	"fmt"
	"log"
)

type IntHeap []int //nolint:recvcheck

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int)) //nolint:forcetypeassert
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

const ErrorEOF = "EOF"

func readInt() (int, error) {
	var value int
	_, err := fmt.Scan(&value)

	if err != nil && err.Error() != ErrorEOF {
		return 0, err //nolint:wrapcheck
	}

	return value, nil
}

func readPreferences(dishAmount int) ([]int, error) {
	preferences := make([]int, dishAmount)
	for i := range dishAmount {
		if _, err := fmt.Scan(&preferences[i]); err != nil {
			if err.Error() != ErrorEOF {
				return nil, err //nolint:wrapcheck
			}

			break
		}
	}

	return preferences, nil
}

func finddishNumberthSmallest(preferences []int, dishNumber int) int {
	h := &IntHeap{} //nolint:varnamelen

	for _, pref := range preferences {
		heap.Push(h, pref)

		if h.Len() > dishNumber {
			heap.Pop(h)
		}
	}

	if h.Len() > 0 {
		return (*h)[0]
	}

	return 0
}

func run() error {
	dishAmount, err := readInt()
	if err != nil {
		return err
	}

	preferences, err := readPreferences(dishAmount)
	if err != nil {
		return err
	}

	dishNumber, err := readInt()
	if err != nil {
		return err
	}

	if dishNumber < 1 || dishNumber > dishAmount {
		return nil
	}

	result := finddishNumberthSmallest(preferences, dishNumber)
	fmt.Println(result)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
