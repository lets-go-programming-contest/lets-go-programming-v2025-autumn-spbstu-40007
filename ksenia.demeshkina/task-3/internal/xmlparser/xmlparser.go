// загрузка и парсинг XML
package xmlparser

import (
	"bytes"
	"encoding/xml" // для разбора XML
	"fmt"          // переводить строку "43,6438" -> 43.6438
	"io"
	"log"     // вывод ошибок log.Panic
	"os"      // для открытия и чтения файла
	"strings" // чтоб заменить запятые на точки

	"golang.org/x/text/encoding/charmap"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

// структура одной валюты
type Valute struct {
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	ValueStr string  `xml:"Value"`
	Value    float64 `xml:"-"` // это поле вычисляемое но не из XML
}

func LoadXML(path string) []Valute {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("Cannot open XML file: %v", err)
	}

	reader := bytes.NewReader(data)

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}

	// парсинг XML
	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil { // разбирает data и заполняет valCurs
		log.Panicf("Cannot parse XML: %v", err)
	}

	// преобразование строки ValueStr в числа Value
	for i := range valCurs.Valutes { // цикл по всем валютам
		valCurs.Valutes[i].Value = parseValue(valCurs.Valutes[i].ValueStr) // для каждой валюты берем ValueStr (строку с запятой), передаем в parseValue и записываем результат в Value
	}

	return valCurs.Valutes // возвращаем готовый срез
}

// перевод строки вроде 43,432 в 43.432
func parseValue(s string) float64 {
	s = strings.Replace(s, ",", ".", 1)
	var val float64
	_, err := fmt.Sscanf(s, "%f", &val)
	if err != nil {
		log.Panicf("Invalid number format in XML: %s", s)
	}

	return val
}
