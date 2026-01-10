package configuration

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	InputFilePath  string `yaml:"input-file"`
	OutputFilePath string `yaml:"output-file"`
}

func LoadApplicationConfig(configurationPath string) (*ApplicationConfig, error) {
	configurationContent, readError := os.ReadFile(configurationPath)
	if readError != nil {
		return nil, fmt.Errorf("reading configuration file at %q: %w", configurationPath, readError)
	}

	applicationConfig := new(ApplicationConfig)

	unmarshalError := yaml.Unmarshal(configurationContent, applicationConfig)
	if unmarshalError != nil {
		return nil, fmt.Errorf("unmarshalling YAML configuration: %w", unmarshalError)
	}

	return applicationConfig, nil
}
