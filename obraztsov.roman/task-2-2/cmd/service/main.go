package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }

func (h IntHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		fmt.Println("Error use int ")
		return
	}
	*h = append(*h, value)
}
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]

	old[n-1] = 0

	old = old[0 : n-1]
	return item

}

func (h IntHeap) Sort() {
	n := h.Len()
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if h.Less(j+1, j) {
				h.Swap(j, j+1)
			}

		}
	}
}

func main() {
	var (
		numberAll int
		//numberOfFood int

	)
	intHeap := &IntHeap{}

	_, err1 := fmt.Scanln(&numberAll)
	if err1 != nil {
		fmt.Println("Invalid Value")
		return
	}

	if !(numberAll >= 1 && numberAll <= 10000) {
		fmt.Println("Number of food out of range")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	line, err3 := reader.ReadString('\n')
	if err3 != nil {
		fmt.Println("Error cant read food values")
		return
	}
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if len(parts) != numberAll {
		fmt.Println("Invalid input")
		return
	}
	f := true

	for _, part := range parts {
		numberOfFood, err := strconv.Atoi(part)
		if err != nil {
			fmt.Println("Invalid value")
			f = false
			return
		}
		if !(numberOfFood >= -10000 && numberOfFood <= 10000) {
			fmt.Println("Invalid value")
			f = false
			return
		}
		heap.Push(intHeap, numberOfFood)

	}
	if f == false {
		return
	}

	intHeap.Sort()
	//reader2:= bufio.NewReader(os.Stdin)
	value, err4 := reader.ReadString('\n')
	if err4 != nil {
		fmt.Println("Error of value food ")
	}
	value = strings.TrimSpace(value)
	valueFood, err5 := strconv.Atoi(value)
	if err5 != nil {
		fmt.Println("Invalid value")
	}

	if !(valueFood <= len(*intHeap) && valueFood > 0) {
		fmt.Println("Invalid value")

	}
	heap := *intHeap
	m := len(heap)

	fmt.Println(heap[m-valueFood])

}
