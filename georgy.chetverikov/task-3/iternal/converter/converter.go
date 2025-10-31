package converter

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/falsefeelings/task-3/iternal/data"
	"gopkg.in/yaml.v3"
)

type Converter struct{}

var ErrUnsupportedOutputFormat = errors.New("unsupported output format")

func New() *Converter {
	return &Converter{}
}

func (c *Converter) Convert(valutes data.Valutes, format string) ([]byte, error) {
	validFormats := map[string]bool{
		"json": true,
		"yaml": true,
		"xml":  true,
	}

	if !validFormats[format] {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOutputFormat, format)
	}

	switch format {
	case "json":
		data, err := json.MarshalIndent(valutes, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal json: %w", err)
		}

		return data, nil
	case "yaml":
		data, err := yaml.Marshal(valutes)
		if err != nil {
			return nil, fmt.Errorf("marshal yaml: %w", err)
		}

		return data, nil
	case "xml":
		type ValCurs struct {
			XMLName xml.Name
			Valutes data.Valutes `xml:"Valute"`
		}

		wrapper := ValCurs{
			XMLName: xml.Name{Local: "ValCurs"},
			Valutes: valutes,
		}

		data, err := xml.MarshalIndent(wrapper, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal xml: %w", err)
		}

		return data, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOutputFormat, format)
	}
}
