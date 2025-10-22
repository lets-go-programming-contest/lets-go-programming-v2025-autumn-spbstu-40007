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

func makeRange() TempRange {
	return TempRange{lowLimit, highLimit}
}

func (r *TempRange) adjust(sign string, t int) bool {
	if sign == ">=" {
		if t > r.low {
			r.low = t
		}
	} else if sign == "<=" {
		if t < r.high {
			r.high = t
		}
	} else {
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
	var sections int
	if _, err := fmt.Scan(&sections); err != nil {
		return
	}

	for s := 0; s < sections; s++ {
		var workers int
		fmt.Scan(&workers)

		limits := makeRange()

		for w := 0; w < workers; w++ {
			var op string
			var val int
			fmt.Scan(&op, &val)

			ok := limits.adjust(op, val)
			if !ok {
				return
			}
			fmt.Println(limits.current())
		}
	}
}
