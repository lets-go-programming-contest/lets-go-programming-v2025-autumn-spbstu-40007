package converter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/treadwave/task-3/internal/structs"
)

func Converter(tempValutes []structs.TempValute) ([]structs.Valute, error) {
	result := make([]structs.Valute, len(tempValutes))

	for index, tempValute := range tempValutes {
		valueString := strings.Replace(tempValute.Value, ",", ".", 1)

		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting string to float64: %w", err)
		}

		result[index] = structs.Valute{
			NumCode:  tempValute.NumCode,
			CharCode: tempValute.CharCode,
			Value:    value,
		}
	}

	return result, nil
}
