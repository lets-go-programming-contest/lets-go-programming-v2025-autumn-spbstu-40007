package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputDir  string `yaml:"output-dir"`
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

type CurrencyOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func loadConfig(configPath string) (Config, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.OutputDir == "" {
		config.OutputDir = "./output"
	}

	return config, nil
}

func parseCurrencies(valCurs ValCurs) ([]CurrencyOutput, error) {
	currencies := make([]CurrencyOutput, 0, len(valCurs.Valutes))
	for _, valute := range valCurs.Valutes {
		numCode := 0
		if valute.NumCode != "" {
			if parsed, err := strconv.Atoi(valute.NumCode); err == nil {
				numCode = parsed
			}
		}
		value, err := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", 1), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for %s: %w", valute.CharCode, err)
		}
		currencies = append(currencies, CurrencyOutput{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}
	return currencies, nil
}

func main() {
	configPath := flag.String("config", "", "Path to the configuration YAML file")
	flag.Parse()

	if *configPath == "" {
		panic("Configuration file path is required via --config flag")
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(config.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to open input file: %v", err))
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unsupported charset: %s", charset)
		}
	}

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse XML: %v", err))
	}

	currencies, err := parseCurrencies(valCurs)
	if err != nil {
		panic(err)
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	err = os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to create output directory: %v", err))
	}

	outputFilePath := filepath.Join(config.OutputDir, config.OutputFile)

	outputData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON: %v", err))
	}

	err = os.WriteFile(outputFilePath, outputData, 0600)
	if err != nil {
		panic(fmt.Sprintf("Failed to write output file: %v", err))
	}

	fmt.Println("Processing completed successfully! Output saved to:", outputFilePath)
}
