package main

import (
	"container/heap"
	"fmt"
	"sort"
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
	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Println("Ошибка ввода N")
		return
	}

	if N < 1 || N > 10000 {
		fmt.Println("Блюда должны быть от 1 до 10000")
		return
	}

	h := &IntHeap{}
	heap.Init(h)

	for i := 0; i < N; i++ {
		var ai int
		_, err = fmt.Scan(&ai)
		if err != nil {
			fmt.Println("Ошибка ввода последовательности")
			return
		}

		if ai < -10000 || ai > 10000 {
			fmt.Println("Последовательность в диапазоне от -10000 до 10000")
			return
		}

		heap.Push(h, ai)
	}

	var k int
	_, err = fmt.Scan(&k)
	if err != nil {
		fmt.Println("Ошибка ввода k")
		return
	}

	if k < 1 || k > N {
		fmt.Printf("K должно быть от 1 до %d\n", N)
		return
	}

	slice := []int(*h)

	sort.Sort(sort.Reverse(sort.IntSlice(slice)))

	result := slice[k-1]
	fmt.Println(slice)
	fmt.Println(result)

}
