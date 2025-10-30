package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppSettings struct {
	SourcePath string `yaml:"input-file"`
	TargetPath string `yaml:"output-file"`
}

func ReadSettings(filepath string) AppSettings {
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Configuration file reading error '%s': %v", filepath, err)
		panic("configuration file access failed: " + err.Error())
	}

	var settings AppSettings
	if unmarshalErr := yaml.Unmarshal(content, &settings); unmarshalErr != nil {
		log.Printf("YAML configuration parsing error: %v", unmarshalErr)
		panic("YAML configuration parsing failed: " + unmarshalErr.Error())
	}

	return settings
}
