package exchangerate

//nolint:gofumpt,gci
import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"task-3/internal/die"
	"task-3/internal/files"

	"golang.org/x/net/html/charset"
)

var errWrapErr = func(err error) error {
	return fmt.Errorf("exchrangerate: %w", err)
}

type RussianFloat float64

func (russianFloat *RussianFloat) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	content := ""
	if err := decoder.DecodeElement(&content, &start); err != nil { //nolint:noinlineerr
		return errWrapErr(err)
	}

	content = strings.ReplaceAll(content, ",", ".")

	result, err := strconv.ParseFloat(content, 64)
	if err != nil {
		return errWrapErr(err)
	}

	*russianFloat = RussianFloat(result)

	return nil
}

type Currency struct {
	NumCode  int          `json:"num_code"  xml:"NumCode"`
	CharCode string       `json:"char_code" xml:"CharCode"`
	Value    RussianFloat `json:"value"     xml:"Value"`
}

type ExchangeRate struct {
	Currencies []Currency `xml:"Valute"`
}

func FromXMLFile(path string) (*ExchangeRate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errWrapErr(err)
	}

	defer func() {
		if err = file.Close(); err != nil { //nolint:noinlineerr
			die.Die(err)
		}
	}()

	result := &ExchangeRate{} //nolint:exhaustruct

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err = decoder.Decode(result); err != nil { //nolint:noinlineerr
		return nil, errWrapErr(err)
	}

	return result, nil
}

func (exchangeRate *ExchangeRate) ToJSONFile(path string) error {
	err := files.CreateIfNotExists(path)
	if err != nil {
		return errWrapErr(err)
	}

	// Magic number? Are you serious?
	file, err := os.OpenFile(path, os.O_WRONLY, 0o600) //nolint:mnd
	if err != nil {
		return errWrapErr(err)
	}

	defer func() {
		if err = file.Close(); err != nil { //nolint:noinlineerr
			die.Die(err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(exchangeRate.Currencies); err != nil { //nolint:noinlineerr
		return errWrapErr(err)
	}

	return nil
}
