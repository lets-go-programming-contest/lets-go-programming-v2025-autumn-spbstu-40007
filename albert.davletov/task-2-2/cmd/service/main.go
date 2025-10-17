package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(scanner *bufio.Scanner) (string, error) {
	scanner.Scan()
	err := scanner.Err()

	return scanner.Text(), err
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
	intX, ok := x.(int)
	if ok {
		*ownHeap = append(*ownHeap, intX)
	}
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
	preferedDishes, err := readInput(reader)

	if err != nil {
		fmt.Println("Error reading input: ", err)

		return
	}

	for _, v := range strings.Fields(preferedDishes) {
		intv, err := strconv.Atoi(v)

		if err != nil {
			fmt.Println("Erorr converting to number: ", err)
		}

		dishesHeap.Push(intv)
	}

	if dishesHeap.Len() != dishesCount {
		fmt.Println("Wrong number of dishes")

		return
	}

	heap.Init(dishesHeap)

	numberOfPreferedDishString, err := readInput(reader)

	if err != nil {
		fmt.Println("Error reading input: ", err)

		return
	}

	numberOfPreferedDish, err := strconv.Atoi(numberOfPreferedDishString)

	if err != nil {
		fmt.Println("Erorr converting to number: ", err)

		return
	}

	for range numberOfPreferedDish - 1 {
		heap.Pop(dishesHeap)
	}

	fmt.Println(heap.Pop(dishesHeap))
}
