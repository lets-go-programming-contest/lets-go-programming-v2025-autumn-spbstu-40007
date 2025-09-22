package main

//nolint:all
import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"

	"task-2-2/internal/functional"
	"task-2-2/internal/scanner"
)

type IntMaxPriorityQueue []int //nolint:all

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
	*priorityQueue = append(*priorityQueue, x.(int)) //nolint:all
}

func (priorityQueue *IntMaxPriorityQueue) Pop() any {
	x := (*priorityQueue)[len(*priorityQueue)-1]
	*priorityQueue = (*priorityQueue)[:len(*priorityQueue)-1]

	return x
}

func main() {
	scanner := scanner.NewScanner()
	scanner.SkipNLines(1)
	as := (IntMaxPriorityQueue)(functional.Map( //nolint:all
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
