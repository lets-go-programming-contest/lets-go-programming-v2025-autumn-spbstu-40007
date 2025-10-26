package exchangerate

//nolint:gofumpt,gci
import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"strconv"
	"strings"

	"task-3/internal/die"
	"task-3/internal/files"

	"golang.org/x/net/html/charset"
)

var (
	errFailedToUnmarshalFloat = errors.New("exchangerate: failed to unmarshal float")
	errFailedToReadXMLFile    = errors.New("exchangerate: failed to read xml file")
	errFailedToWriteJSONFile  = errors.New("exchangerate: failed to write json file")
)

type RussianFloat float64

func (russianFloat *RussianFloat) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	content := ""
	if err := decoder.DecodeElement(&content, &start); err != nil {
		return errFailedToUnmarshalFloat
	}

	content = strings.ReplaceAll(content, ",", ".")

	result, err := strconv.ParseFloat(content, 64)
	if err != nil {
		return errFailedToUnmarshalFloat
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
		return nil, fmt.Errorf("%w: %w", errFailedToReadXMLFile, err)
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
		return nil, fmt.Errorf("%w: %w", errFailedToReadXMLFile, err)
	}

	return result, nil
}

func (exchangeRate *ExchangeRate) ToJSONFile(path string) error {
	err := files.CreateIfNotExists(path)
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToWriteJSONFile, err)
	}

	// Magic number? Are you serious?
	file, err := os.OpenFile(path, os.O_WRONLY, 0o600) //nolint:mnd
	if err != nil {
		return fmt.Errorf("%w: %w", errFailedToWriteJSONFile, err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			die.Die(err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(exchangeRate.Currencies); err != nil {
		return fmt.Errorf("%w: %w", errFailedToWriteJSONFile, err)
	}

	return nil
}
