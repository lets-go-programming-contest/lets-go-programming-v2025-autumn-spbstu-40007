package vp

import (
	"sort"

	"github.com/tntkatz/task-3/internal/data"
)

func ProcessValute(valCurs data.ValCurs) ([]data.Valute, error) {
	sortedValutes := make(data.Valutes, len(valCurs.Valutes))

	copy(sortedValutes, valCurs.Valutes)

	sort.Sort(sortedValutes)

	return sortedValutes, nil
}
