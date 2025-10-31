package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppSettings struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ConfigurationManager struct{}

func LoadSettings(filePath string) AppSettings {
	fileContent, readErr := readFileContent(filePath)
	if readErr != nil {
		handleConfigError(filePath, readErr)
	}

	settings, parseErr := parseYAMLConfig(fileContent)
	if parseErr != nil {
		handleParseError(parseErr)
	}

	return settings
}

func readFileContent(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("невозможно прочитать файл '%s': %w", filePath, err)
	}

	return content, nil
}

func parseYAMLConfig(content []byte) (AppSettings, error) {
	var settings AppSettings
	err := yaml.Unmarshal(content, &settings)

	if err != nil {
		return AppSettings{}, fmt.Errorf("ошибка разбора YAML: %w", err)
	}

	return settings, nil
}

func handleConfigError(filePath string, err error) {
	errorMsg := fmt.Sprintf(
		"Ошибка чтения конфигурационного файла '%s': %v",
		filePath,
		err,
	)
	fmt.Println(errorMsg)
	panic(fmt.Errorf("ошибка загрузки конфигурации: %w", err))
}

func handleParseError(err error) {
	fmt.Printf("Ошибка декодирования YAML конфигурации: %v\n", err)
	panic(fmt.Errorf("ошибка обработки конфигурации: %w", err))
}
