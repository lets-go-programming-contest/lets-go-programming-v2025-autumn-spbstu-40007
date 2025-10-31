package output

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/itsdasha/task-3/internal/data"
	"gopkg.in/yaml.v3"
)

var ErrUnsupportedFormat = errors.New("unsupported format")

const dirPerm = 0o755
const filePerm = 0o600

func Save(results []data.OutputCurrency, path, format string) {
	var content []byte
	var err error

	switch strings.ToLower(format) {
	case "json":
		content, err = json.MarshalIndent(results, "", "  ")
	case "yaml":
		content, err = yaml.Marshal(results)
	case "xml":
		type wrapper struct {
			XMLName xml.Name              `xml:"ValCurs"`
			Items   []data.OutputCurrency `xml:"Valute"`
		}
		var w wrapper
		w.XMLName = xml.Name{Space: "", Local: "ValCurs"}
		w.Items = results
		content, err = xml.MarshalIndent(w, "", "  ")
		content = []byte(xml.Header + string(content))
	default:
		panic(fmt.Errorf("%w: %s", ErrUnsupportedFormat, format))
	}

	if err != nil {
		panic(fmt.Errorf("encoding error: %w", err))
	}

	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		panic(fmt.Errorf("cannot create directory: %w", err))
	}

	if err := os.WriteFile(path, content, filePerm); err != nil {
		panic(fmt.Errorf("cannot write file: %w", err))
	}
}
