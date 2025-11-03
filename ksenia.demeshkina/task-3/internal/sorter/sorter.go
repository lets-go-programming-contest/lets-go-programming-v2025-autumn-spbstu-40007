package sorter

import (
	"sort"

	"github.com/ksuah/task-3/internal/xmlparser"
)

func SortByValueDesc(valutes []xmlparser.Valute) []xmlparser.Valute {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})

	return valutes
}
