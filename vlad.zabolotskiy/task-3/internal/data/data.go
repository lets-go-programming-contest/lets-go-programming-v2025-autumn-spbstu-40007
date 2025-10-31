package data

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Currency struct {
	NumCode  int     `json:"num_code" xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float32 `json:"value" xml:"Value"`
}

type ValCurs struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}

func LoadFromXML(path string) (*ValCurs, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	xmlContent := strings.ReplaceAll(string(data), ",", ".")

	var valCurs ValCurs
	err = xml.Unmarshal([]byte(xmlContent), &valCurs)

	if err != nil {
		return nil, err
	}

	return &valCurs, nil
}

func (vc *ValCurs) SortByValue() {
	compareCurrencies := func(i, j int) bool {
		valueI := vc.Currencies[i].Value
		valueJ := vc.Currencies[j].Value

		if valueI > valueJ {
			return true
		}

		return false
	}

	sort.Slice(vc.Currencies, compareCurrencies)
}

func (vc *ValCurs) SaveToJSON(path string) error {
	folderPath := filepath.Dir(path)
	os.MkdirAll(folderPath, 0755)

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	jsonData, err := json.MarshalIndent(vc.Currencies, "", "  ")

	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	return err
}
