package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type OutputCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

var configPath = flag.String("config", "config.yaml", "Path to configuration file")

func main() {
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}

	valCurs, err := parseXML(config.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Error parsing XML: %v", err))
	}

	outputCurrencies := convertAndSort(valCurs)

	err = saveJSON(config.OutputFile, outputCurrencies)
	if err != nil {
		panic(fmt.Sprintf("Error saving JSON: %v", err))
	}

	fmt.Println("Successfully processed and saved currencies.")
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "input-file:") {
			config.InputFile = strings.Trim(strings.TrimPrefix(line, "input-file:"), " \"")
		} else if strings.HasPrefix(line, "output-file:") {
			config.OutputFile = strings.Trim(strings.TrimPrefix(line, "output-file:"), " \"")
		}
	}

	if config.InputFile == "" || config.OutputFile == "" {
		return nil, fmt.Errorf("invalid config format")
	}

	return &config, nil
}

func parseXML(filePath string) (*ValCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	var valCurs ValCurs
	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &valCurs, nil
}

func convertAndSort(valCurs *ValCurs) []OutputCurrency {
	var output []OutputCurrency

	for _, v := range valCurs.Valutes {
		numCode, err := strconv.Atoi(v.NumCode)
		if err != nil {
			continue
		}

		valueStr := strings.Replace(v.Value, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		output = append(output, OutputCurrency{
			NumCode:  numCode,
			CharCode: v.CharCode,
			Value:    value,
		})
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Value > output[j].Value
	})

	return output
}

func saveJSON(filePath string, data []OutputCurrency) error {
	dir := filePath[:strings.LastIndex(filePath, "/")]
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
