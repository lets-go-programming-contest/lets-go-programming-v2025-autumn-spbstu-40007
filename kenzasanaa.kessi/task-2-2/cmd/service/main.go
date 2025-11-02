package main

import (
	"container/heap"
	"errors"
	"fmt"

	"kenzasanaa.kessi/task-2-2/internal/intheap"
)

var (
	ErrInvalidPosition  = errors.New("invalid position")
	ErrHeapEmpty        = errors.New("heap became empty unexpectedly")
	ErrNoResult         = errors.New("unable to retrieve result from heap")
	ErrInvalidDataType  = errors.New("invalid data type in heap")
	ErrInvalidDishCount = errors.New("invalid dish count")
	ErrInvalidKValue    = errors.New("invalid k value")
)

func getKthMaximum(values []int, position int) (int, error) {
	if position <= 0 || position > len(values) {
		return 0, fmt.Errorf("%w: position %d for slice size %d", ErrInvalidPosition, position, len(values))
	}

	maxHeap := &intheap.CustomHeap{}
	heap.Init(maxHeap)

	for _, val := range values {
		heap.Push(maxHeap, val)
	}

	for range position - 1 {
		if maxHeap.Size() == 0 {
			return 0, ErrHeapEmpty
		}

		heap.Pop(maxHeap)
	}

	result := heap.Pop(maxHeap)
	if result == nil {
		return 0, ErrNoResult
	}

	finalResult, valid := result.(int)
	if !valid {
		return 0, ErrInvalidDataType
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
		fmt.Println(ErrInvalidDishCount)

		return
	}

	ratings := make([]int, totalDishes)
	for i := range totalDishes {
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
		fmt.Printf("Error: %v - position %d for %d dishes\n", ErrInvalidKValue, targetPosition, totalDishes)

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
