package main

import (
	"container/heap"
	"fmt"

	"kenzasanaa.kessi/task-2-2/internal/intheap"
)

func getKthMaximum(values []int, position int) (int, error) {
	if position <= 0 || position > len(values) {
		return 0, fmt.Errorf("position %d is invalid for slice of size %d", position, len(values))
	}

	maxHeap := &intheap.CustomHeap{}
	heap.Init(maxHeap)

	for _, val := range values {
		heap.Push(maxHeap, val)
	}

	for i := 0; i < position-1; i++ {
		if maxHeap.Size() == 0 {
			return 0, fmt.Errorf("heap became empty unexpectedly")
		}
		heap.Pop(maxHeap)
	}

	result := heap.Pop(maxHeap)
	if result == nil {
		return 0, fmt.Errorf("unable to retrieve result from heap")
	}

	finalResult, valid := result.(int)
	if !valid {
		return 0, fmt.Errorf("invalid data type in heap")
	}

	return finalResult, nil
}

func executeProgram() {
	var totalDishes int
	_, scanErr := fmt.Scan(&totalDishes)
	if scanErr != nil {
		fmt.Printf("Error reading dish count: %v\n", scanErr)
		return
	}

	if totalDishes <= 0 {
		fmt.Println("No dishes available")
		return
	}


	ratings := make([]int, totalDishes)
	for i := 0; i < totalDishes; i++ {
		_, ratingErr := fmt.Scan(&ratings[i])
		if ratingErr != nil {
			fmt.Printf("Error reading rating %d: %v\n", i+1, ratingErr)
			return
		}
	}

	var targetPosition int
	_, posErr := fmt.Scan(&targetPosition)
	if posErr != nil {
		fmt.Printf("Error reading target position: %v\n", posErr)
		return
	}


	if targetPosition > totalDishes || targetPosition <= 0 {
		fmt.Printf("Position %d is not valid for %d dishes\n", targetPosition, totalDishes)
		return
	}

	kthLargest, calcErr := getKthMaximum(ratings, targetPosition)
	if calcErr != nil {
		fmt.Printf("Calculation error: %v\n", calcErr)
		return
	}

	fmt.Println(kthLargest)
}

func main() {
	defer func() {
		if recovery := recover(); recovery != nil {
			fmt.Printf("Program recovered from panic: %v\n", recovery)
		}
	}()

	executeProgram()
}
