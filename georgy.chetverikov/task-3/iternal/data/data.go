package data

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Valutes []Valute

type Valute struct {
	ID       string      `json:"id"        xml:"ID,attr"`
	NumCode  int         `json:"num_code"  xml:"NumCode"`
	CharCode string      `json:"char_code" xml:"CharCode"`
	Nominal  int         `json:"nominal"   xml:"Nominal"`
	Name     string      `json:"name"      xml:"Name"`
	Value    customFloat `json:"value"     xml:"Value"`
}

type customFloat float64

var ErrUnsupportedCharset = errors.New("unsupported charset")

func (f *customFloat) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var XMLdata string

	if err := decoder.DecodeElement(&XMLdata, &start); err != nil {
		return fmt.Errorf("decoding xml-file: %w", err)
	}

	decodedData := strings.Replace(XMLdata, ",", ".", 1)

	parsedData, err := strconv.ParseFloat(decodedData, 64)
	if err != nil {
		return fmt.Errorf("parsing data: %w", err)
	}

	*f = customFloat(parsedData)

	return nil
}

func (v Valutes) Len() int           { return len(v) }
func (v Valutes) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Valutes) Less(i, j int) bool { return v[i].Value > v[j].Value }

func ParseXML(data []byte) (Valutes, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, charset)
	}

	type tempValCourse struct {
		Date    string   `xml:"Date,attr"`
		Name    string   `xml:"name,attr"`
		Valutes []Valute `xml:"Valute"`
	}

	var temp tempValCourse

	err := decoder.Decode(&temp)
	if err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	valutes := temp.Valutes
	sort.Sort(Valutes(valutes))

	return valutes, nil
}
