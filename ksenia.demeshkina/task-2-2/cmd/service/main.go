package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var count int
	_, err := fmt.Scan(&count)

	if err != nil {
		fmt.Println("Ошибка ввода count")

		return
	}

	if count < 1 || count > 10000 {
		fmt.Println("Блюда должны быть от 1 до 10000")

		return
	}

	dishesHeap := &IntHeap{}
	heap.Init(dishesHeap)

	for range count {
		var dishes int
		_, err = fmt.Scan(&dishes)

		if err != nil {
			fmt.Println("Ошибка ввода последовательности")

			return
		}

		if dishes < -10000 || dishes > 10000 {
			fmt.Println("Последовательность в диапазоне от -10000 до 10000")

			return
		}

		heap.Push(dishesHeap, dishes)
	}

	var favorite int
	_, err = fmt.Scan(&favorite)

	if err != nil {
		fmt.Println("Ошибка ввода favorite")

		return
	}

	if favorite < 1 || favorite > count {
		fmt.Printf("favorite должно быть от 1 до %d\n", count)

		return
	}

	slice := []int(*dishesHeap)

	sort.Sort(sort.Reverse(sort.IntSlice(slice)))

	result := slice[favorite-1]
	fmt.Println(result)
}
