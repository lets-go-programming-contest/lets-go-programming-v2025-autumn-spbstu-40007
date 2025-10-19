package vp

import (
	"sort"

	"github.com/tntkatz/task-3/internal/data"
)

func ValuteProcess(valCurs data.ValCurs) ([]data.Valute, error) {
	sortedValutes := make([]data.Valute, len(valCurs.Valute))

	copy(sortedValutes, valCurs.Valute)

	sort.Sort(data.ByValue(sortedValutes))

	return sortedValutes, nil
}
