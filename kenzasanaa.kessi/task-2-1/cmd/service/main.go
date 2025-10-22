package main

import (
	"fmt"
)

const (
	lowLimit  = 15
	highLimit = 30
)

type TempRange struct {
	low  int
	high int
}

func newTempRange() TempRange {
	return TempRange{low: lowLimit, high: highLimit}
}

func (r *TempRange) adjust(operator string, tempValue int) bool {
	switch operator {
	case ">=":
		if tempValue > r.low {
			r.low = tempValue
		}
	case "<=":
		if tempValue < r.high {
			r.high = tempValue
		}
	default:
		return false
	}

	return true
}

func (r *TempRange) current() int {
	if r.low <= r.high {
		return r.low
	}

	return -1
}

func main() {
	var sectionCount int
	if _, err := fmt.Scan(&sectionCount); err != nil {
		return
	}

	for range make([]struct{}, sectionCount) {
		var workerCount int
		if _, err := fmt.Scan(&workerCount); err != nil {
			return
		}

		limits := newTempRange()

		for range make([]struct{}, workerCount) {
			var (
				op  string
				val int
			)
			if _, err := fmt.Scan(&op, &val); err != nil {
				return
			}

			if !limits.adjust(op, val) {
				return
			}

			fmt.Println(limits.current())
		}
	}
}
