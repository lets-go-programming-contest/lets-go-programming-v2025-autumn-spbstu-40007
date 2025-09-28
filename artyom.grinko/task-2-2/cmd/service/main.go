package main

// Newline for separation std packages from others seems crazy to linter.
//nolint:gofumpt
import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"

	"task-2-2/internal/functional"
	"task-2-2/internal/scanner"
)

// Linter doesn't really like code snippet from https://pkg.go.dev/container/heap.
type IntMaxPriorityQueue []int //nolint:recvcheck

func (priorityQueue IntMaxPriorityQueue) Len() int {
	return len(priorityQueue)
}

func (priorityQueue IntMaxPriorityQueue) Less(i, j int) bool {
	return priorityQueue[i] > priorityQueue[j]
}

func (priorityQueue IntMaxPriorityQueue) Swap(i, j int) {
	priorityQueue[i], priorityQueue[j] = priorityQueue[j], priorityQueue[i]
}

func (priorityQueue *IntMaxPriorityQueue) Push(x any) {
	// Type checking where only one type used is redundant.
	*priorityQueue = append(*priorityQueue, x.(int)) //nolint:forcetypeassert
}

func (priorityQueue *IntMaxPriorityQueue) Pop() any {
	x := (*priorityQueue)[len(*priorityQueue)-1]
	*priorityQueue = (*priorityQueue)[:len(*priorityQueue)-1]

	return x
}

func main() {
	scanner := scanner.NewScanner()
	scanner.SkipNLines(1)
	// Variable has same name as in task.
	as := (IntMaxPriorityQueue)(functional.Map( //nolint:varnamelen
		strings.Fields(scanner.Read()),
		func(x string) int {
			y, _ := strconv.Atoi(x)

			return y
		}),
	)
	heap.Init(&as)

	k, _ := strconv.Atoi(scanner.Read())
	for range k - 1 {
		heap.Pop(&as)
	}

	fmt.Println(heap.Pop(&as))
}
