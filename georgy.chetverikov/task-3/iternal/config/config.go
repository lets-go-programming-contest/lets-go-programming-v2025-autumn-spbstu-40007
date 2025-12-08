package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile    string `yaml:"input-file"`
	OutputFile   string `yaml:"output-file"`
	OutputFormat string `yaml:"output-format"`
}

var (
	ErrInputFileRequired   = errors.New("input-file is required")
	ErrOutputFileRequired  = errors.New("output-file is required")
	ErrInvalidInputFile    = errors.New("input file must be a XML file")
	ErrInvalidOutputFormat = errors.New("output-format must be one of: json, yaml, xml")
)

func Read(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config %q: %w", configPath, err)
	}

	config := new(Config)

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("unable to unmarshall data: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid YAML-config: %w", err)
	}

	if config.OutputFormat == "" {
		config.OutputFormat = "json"
	}

	return config, nil
}

func validateConfig(config *Config) error {
	if config.InputFile == "" {
		return ErrInputFileRequired
	}

	if config.OutputFile == "" {
		return ErrOutputFileRequired
	}

	if !isXMLFile(config.InputFile) {
		return fmt.Errorf("%w, got: %s", ErrInvalidInputFile, config.InputFile)
	}

	validFormats := map[string]bool{
		"json": true,
		"yaml": true,
		"xml":  true,
	}

	if config.OutputFormat != "" && !validFormats[config.OutputFormat] {
		return fmt.Errorf("%w, got: %s", ErrInvalidOutputFormat, config.OutputFormat)
	}

	return nil
}

func isXMLFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	return ext == ".xml"
}
