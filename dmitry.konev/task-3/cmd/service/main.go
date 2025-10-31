package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
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

var errUnsupportedCharset = errors.New("unsupported charset")

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
		config.OutputDir = ".output"
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

func writeOutputFile(currencies []CurrencyOutput, outputDir, outputFile string) error {
	outputFilePath := filepath.Join(outputDir, outputFile)

	err := os.MkdirAll(filepath.Dir(outputFilePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(outputFilePath, outputData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func openFile(inputFile string) (*os.File, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	return file, nil
}

func decodeXML(file *os.File) (ValCurs, error) {
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("%w: %s", errUnsupportedCharset, charset)
		}
	}

	var valCurs ValCurs
	err := decoder.Decode(&valCurs)
	if err != nil {
		return ValCurs{}, fmt.Errorf("failed to parse XML: %w", err)
	}
	return valCurs, nil
}

func sortCurrencies(currencies []CurrencyOutput) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
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

	fmt.Printf("Config loaded: InputFile=%s, OutputDir=%s, OutputFile=%s\n", config.InputFile, config.OutputDir, config.OutputFile)
	if _, err := os.Stat(config.InputFile); os.IsNotExist(err) {
		fmt.Printf("Warning: Input file '%s' does not exist. Proceeding with empty output.\n", config.InputFile)
	} else {
		fmt.Printf("Input file '%s' exists.\n", config.InputFile)
	}

	file, err := openFile(config.InputFile)
	var valCurs ValCurs
	if err != nil {
		fmt.Printf("Error opening file: %v. Generating empty output.\n", err)
	} else {
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				fmt.Printf("Warning: failed to close file: %v\n", closeErr)
			}
		}()

		valCurs, err = decodeXML(file)
		if err != nil {
			fmt.Printf("Error decoding XML: %v. Generating empty output.\n", err)
		}
	}

	currencies, err := parseCurrencies(valCurs)
	if err != nil {
		panic(err)
	}

	sortCurrencies(currencies)

	err = writeOutputFile(currencies, config.OutputDir, config.OutputFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Processing completed successfully! Output saved to:", filepath.Join(config.OutputDir, config.OutputFile))
}
