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

func (c *Converter) Convert(valCourse *data.ValCourse, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(valCourse, "", " ")
	case "yaml":
		return yaml.Marshal(valCourse)
	case "xml":
		return xml.MarshalIndent(valCourse, "", " ")
	default:
		return nil, fmt.Errorf("Unsupported output format: %s", format)
	}
}
