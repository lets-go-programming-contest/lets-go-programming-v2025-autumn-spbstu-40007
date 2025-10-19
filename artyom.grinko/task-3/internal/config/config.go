package config

//nolint:gofumpt,gci
import (
	"errors"
	"fmt"
	"os"

	"task-3/internal/die"

	yaml "github.com/goccy/go-yaml"
)

var errDidNotFindExpectedKey = errors.New("config: did not find expected key")

func wrapErr(err error) error {
	return fmt.Errorf("config: %w", err)
}

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, wrapErr(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			die.Die(err)
		}
	}()

	result := &Config{} //nolint:exhaustruct
	decoder := yaml.NewDecoder(file)

	if err = decoder.Decode(result); err != nil {
		return nil, errDidNotFindExpectedKey
	}

	return result, nil
}
