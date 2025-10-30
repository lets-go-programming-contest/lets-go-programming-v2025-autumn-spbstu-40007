package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/falsefeelings/task-3/iternal/data"
	"gopkg.in/yaml.v3"
)

type Converter struct{}

func New() *Converter {
	return &Converter{}
}

func (c *Converter) Convert(valutes data.Valutes, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(valutes, "", "  ")
	case "yaml":
		return yaml.Marshal(valutes)
	case "xml":
		type ValCurs struct {
			XMLName xml.Name     `xml:"ValCurs"`
			Valutes data.Valutes `xml:"Valute"`
		}
		wrapper := ValCurs{Valutes: valutes}
		return xml.MarshalIndent(wrapper, "", "  ")
	default:
		return nil, fmt.Errorf("Unsupported output format: %s", format)
	}
}
