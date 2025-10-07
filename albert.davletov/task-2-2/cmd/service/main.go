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

func (ownHeap *IntHeap) Len() int {
	return len(*ownHeap)
}

func (ownHeap *IntHeap) Less(i, j int) bool {
	return (*ownHeap)[i] > (*ownHeap)[j]
}

func (ownHeap *IntHeap) Swap(i, j int) {
	(*ownHeap)[i], (*ownHeap)[j] = (*ownHeap)[j], (*ownHeap)[i]
}

func (ownHeap *IntHeap) Push(x interface{}) {
	intX, err := x.(int)
	if !err {
		return
	}

	*ownHeap = append(*ownHeap, intX)
}

func (ownHeap *IntHeap) Pop() interface{} {
	old := *ownHeap
	n := len(old)
	x := old[n-1]
	*ownHeap = old[0 : n-1]

	return x
}

func main() {
	dishesHeap := &IntHeap{}

	var dishesCount int
	_, err := fmt.Scanln(&dishesCount)

	if err != nil {
		fmt.Println(err)

		return
	}

	reader := bufio.NewScanner(os.Stdin)
	preferedDishes := readInput(reader)

	for _, v := range strings.Fields(preferedDishes) {
		intv, _ := strconv.Atoi(v)
		dishesHeap.Push(intv)
	}

	if dishesHeap.Len() != dishesCount {
		fmt.Println("Wrong number of dishes")

		return
	}

	heap.Init(dishesHeap)

	numberOfPreferedDish, _ := strconv.Atoi(readInput(reader))
	for range numberOfPreferedDish - 1 {
		heap.Pop(dishesHeap)
	}

	out := heap.Pop(dishesHeap)
	fmt.Println(out)
}
