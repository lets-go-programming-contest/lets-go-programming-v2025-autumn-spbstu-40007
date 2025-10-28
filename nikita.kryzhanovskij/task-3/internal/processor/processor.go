package processor

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"nikita.kryzhanovskij/task-3/internal/models"
)

func Process(valCurs *models.ValCurs) ([]models.ValuteOutput, error) {
	results := make([]models.ValuteOutput, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		value, err := parseValue(valute.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value for %s: %w", valute.CharCode, err)
		}

		results = append(results, models.ValuteOutput{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	return results, nil
}

func parseValue(value string) (float64, error) {
	value = strings.ReplaceAll(value, ",", ".")
	return strconv.ParseFloat(value, 64)
}
