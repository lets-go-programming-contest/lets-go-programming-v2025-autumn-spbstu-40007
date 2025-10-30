package config

import (
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

func FindAndRead() (*Config, error) {

	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("getting executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)

	configFile, err := findYAMLFiles(exeDir)
	if err != nil {
		return nil, err
	}

	if configFile == "" {
		return nil, fmt.Errorf("There are no YAML files in the current direcroty")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read input file %q: %w", configFile, err)
	}

	config := new(Config)

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("Unable to unmarshall data: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid YAML-config: %w", err)
	}

	if config.OutputFormat == "" {
		config.OutputFormat = "json"
	}

	return config, nil
}

func findYAMLFiles(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("reading directory: %w", path, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasSuffix(strings.ToLower(name), ".yaml") ||
				strings.HasSuffix(strings.ToLower(name), ".yml") {

				return filepath.Join(path, name), nil
			}
		}
	}

	return "", fmt.Errorf("There are no YAML files in the directory %q", path)
}

func validateConfig(config *Config) error {

	if config.InputFile == "" {
		return fmt.Errorf("input-file is required")
	}
	if config.OutputFile == "" {
		return fmt.Errorf("output-file is required")
	}

	if !isXMLFile(config.InputFile) {
		return fmt.Errorf("input file must be a XML file, got: %s", config.InputFile)
	}

	validFormats := map[string]bool{
		"json": true,
		"yaml": true,
		"xml":  true,
	}
	if config.OutputFormat != "" && !validFormats[config.OutputFormat] {
		return fmt.Errorf("output-format must be one of: json, yaml,xml, got: %s", config.OutputFormat)
	}

	return nil
}

func isXMLFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".xml"
}
