package main

import (
	"container/heap"
	"fmt"
)

type Heap struct {
	numbers []int
}

func (heap Heap) Len() int {
	return len(heap.numbers)
}

func (heap Heap) Less(i, j int) bool {
	return heap.numbers[i] > heap.numbers[j]
}

func (heap Heap) Swap(i, j int) {
	heap.numbers[i], heap.numbers[j] = heap.numbers[j], heap.numbers[i]
}

func (heap *Heap) Push(x any) {
	heap.numbers = append(heap.numbers, x.(int))
}

func (heap *Heap) Pop() any {
	x := heap.numbers[len(heap.numbers)-1]
	heap.numbres = heap.numbres[:len(heap.numbers)-1]

	return x
}

func (heap Heap) Peek() int {
	return heap.numbers[0]
}

func main() {
	var numbersAmount, targetPosition int

	if _, err := fmt.Scan(&numbersAmount); err != nil {
		fmt.Println("Error: couldn't read the amount of numbers. Numbers must be integers only")
		return
	}

	numbersSequence := make([]int, numbersAmount)

	for index := range numbersSequence {
		if _, err := fmt.Scan(&numbersSequence[index]); err != nil {
			fmt.Println("Error: numbers in the sequence must be integer only")
			return
		}
	}

	if _, err := fmt.Scan(&targetPosition); err != nil {
		fmt.Println("Error: number you are trying to enter must be integer only")
		return
	}

	numbersHeap := &Heap{numbers: []int{}}
	heap.Init(numbersHeap)

	for _, value := range numbersSequence {
		heap.Push(numbersHeap, value)
	}

	for index := 1; index < targetPosition; index++ {
		heap.Pop(numbersHeap)
	}

	fmt.Println(numbersHeap.Peek())
}
