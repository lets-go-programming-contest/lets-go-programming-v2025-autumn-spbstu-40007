package sort

import (
	"sort"

	"github.com/treadwave/task-3/internal/structs"
)

func Sort(valutes []structs.Valute) []structs.Valute {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})

	return valutes
}
