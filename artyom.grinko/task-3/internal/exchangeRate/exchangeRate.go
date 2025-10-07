package exchangeRate

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"strconv"
	"strings"

	"task-3/internal/die"
	"task-3/internal/file"

	"golang.org/x/net/html/charset"
)

type RussianFloat float64

func (russianFloat *RussianFloat) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	content := ""
	if err := decoder.DecodeElement(&content, &start); err != nil {
		return err
	}

	content = strings.ReplaceAll(content, ",", ".")
	result, err := strconv.ParseFloat(content, 64)
	if err != nil {
		return err
	}

	*russianFloat = RussianFloat(result)

	return nil
}

type Currency struct {
	NumCode  int          `xml:"NumCode" json:"num_code"`
	CharCode string       `xml:"CharCode" json:"char_code"`
	Value    RussianFloat `xml:"Value" json:"value"`
}

type ExchangeRate struct {
	Currencies []Currency `xml:"Valute"`
}

func FromXMLFile(path string) (*ExchangeRate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			die.Die(err)
		}
	}()

	result := &ExchangeRate{} //nolint:exhaustruct
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (exchangeRate *ExchangeRate) ToJSONFile(path string) error {
	err := file.CreateIfNotExists(path)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			die.Die(err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(exchangeRate.Currencies); err != nil {
		return err
	}

	return nil
}
