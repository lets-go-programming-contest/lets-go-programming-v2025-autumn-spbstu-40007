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

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

var ErrUnsupportedFormat = errors.New("unsupported format")

type xmlWrapper struct {
	XMLName xml.Name              `xml:"ValCurs"`
	Items   []data.OutputCurrency `xml:"Valute"`
}

func Save(results []data.OutputCurrency, path, format string) {
	var content []byte
	var err error

	switch strings.ToLower(format) {
	case "json":
		content, err = json.MarshalIndent(results, "", "  ")
	case "yaml":
		content, err = yaml.Marshal(results)
	case "xml":
		wrapper := xmlWrapper{
			XMLName: xml.Name{Local: "ValCurs"},
			Items:   results,
		}
		content, err = xml.MarshalIndent(wrapper, "", "  ")
		if err == nil {
			content = []byte(xml.Header + string(content))
		}
	default:
		panic(fmt.Errorf("%w: %s", ErrUnsupportedFormat, format))
	}

	if err != nil {
		panic(fmt.Errorf("encoding error: %w", err))
	}

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		err = os.MkdirAll(dir, dirPerm)
		if err != nil {
			panic(fmt.Errorf("cannot create directory: %w", err))
		}
	}

	err = os.WriteFile(path, content, filePerm)
	if err != nil {
		panic(fmt.Errorf("cannot write file: %w", err))
	}
}
