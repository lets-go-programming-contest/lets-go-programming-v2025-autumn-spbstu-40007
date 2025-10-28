// сортировка валют
package sorter

import (
	"sort"

	"myapp/xmlparser"
)

func SortByValueDesc(valutes []xmlparser.Valute) []xmlparser.Valute {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})

	return valutes
}
