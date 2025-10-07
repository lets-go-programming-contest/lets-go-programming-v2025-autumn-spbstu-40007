package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(scanner *bufio.Scanner) string {
	scanner.Scan()

	err := scanner.Err()
	if err != nil {
		fmt.Println("Error reading input")
		os.Exit(0)
	}

	return scanner.Text()
}

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	h := &IntHeap{}

	var dishesCount int
	_, err := fmt.Scanln(&dishesCount)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewScanner(os.Stdin)
	preferedDishes := readInput(reader)

	for _, v := range strings.Fields(preferedDishes) {
		intV, _ := strconv.Atoi(v)
		h.Push(intV)
	}

	if h.Len() != dishesCount {
		fmt.Println("Wrong number of dishes")
	}

	numberOfPreferedDish, _ := strconv.Atoi(readInput(reader))
	heap.Init(h)

	for range numberOfPreferedDish - 1 {
		heap.Pop(h)
	}

	out := heap.Pop(h)
	fmt.Println(out)
}
