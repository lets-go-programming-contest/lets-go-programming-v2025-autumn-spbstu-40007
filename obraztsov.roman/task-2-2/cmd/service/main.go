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

func (h *IntHeap) Len() int { return len(*h) }

func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }

func (h *IntHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)

	if !ok {
		fmt.Println("Error use int")

		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = 0
	*h = old[0 : n-1]

	return item
}

func (h *IntHeap) Sort() {
	n := h.Len()
	for i := range n - 1 {
		for j := range n - 1 - i {
			if h.Less(j+1, j) {
				h.Swap(j, j+1)
			}
		}
	}
}

func readNumberAll() (int, bool) {
	var numberAll int

	_, err1 := fmt.Scanln(&numberAll)

	if err1 != nil {
		fmt.Println("Invalid Value")

		return 0, false
	}

	return numberAll, true
}

func readFoodValues(reader *bufio.Reader, numberAll int) ([]string, bool) {

	line, err3 := reader.ReadString('\n')

	if err3 != nil {
		fmt.Println("Error cant read food values")

		return nil, false
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")

	if len(parts) != numberAll {
		fmt.Println("Invalid input")

		return nil, false
	}

	return parts, true
}

func processFood(parts []string, intHeap *IntHeap) bool {
	for _, part := range parts {
		numberOfFood, err := strconv.Atoi(part)

		if err != nil {
			fmt.Println("Invalid value")

			return false
		}

		if numberOfFood < -10000 || numberOfFood > 10000 {
			fmt.Println("Invalid value")

			return false
		}

		heap.Push(intHeap, numberOfFood)
	}

	intHeap.Sort()

	return true
}

func readEndProcessFood(reader *bufio.Reader, intHeap *IntHeap) (int, bool) {
	value, err4 := reader.ReadString('\n')

	if err4 != nil {
		fmt.Println("Error of value food")

		return 0, false
	}

	value = strings.TrimSpace(value)

	valueFood, err5 := strconv.Atoi(value)

	if err5 != nil {
		fmt.Println("Invalid value")

		return 0, false
	}

	if valueFood <= 0 || valueFood > len(*intHeap) {
		fmt.Println("Invalid value")

		return 0, false
	}

	return valueFood, true
}

func main() {

	intHeap := &IntHeap{}

	numberAll, ok := readNumberAll()

	if !ok {

		return
	}

	reader := bufio.NewReader(os.Stdin)

	parts, ok2 := readFoodValues(reader, numberAll)

	if !ok2 {

		return
	}

	ok3 := processFood(parts, intHeap)

	if !ok3 {

		return
	}

	valueFood, ok4 := readEndProcessFood(reader, intHeap)

	if !ok4 {

		return
	}

	heapCopy := *intHeap
	length := len(heapCopy)

	fmt.Println(heapCopy[length-valueFood])
}
