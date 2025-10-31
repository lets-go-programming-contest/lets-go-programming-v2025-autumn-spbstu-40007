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

func main() {
	configPath := flag.String("config", "", "Path to the configuration YAML file")
	flag.Parse()

	if *configPath == "" {
		panic("Configuration file path is required via --config flag")
	}

	configData, err := os.ReadFile(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file: %v", err))
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse config file: %v", err))
	}

	file, err := os.Open(config.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to open input file: %v", err))
	}
	defer file.Close()

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

	var currencies []CurrencyOutput
	for _, valute := range valCurs.Valutes {
		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			panic(fmt.Sprintf("Invalid NumCode: %v", err))
		}
		value, err := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", 1), 64)
		if err != nil {
			panic(fmt.Sprintf("Invalid Value: %v", err))
		}
		currencies = append(currencies, CurrencyOutput{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	outputDir := filepath.Dir(config.OutputFile)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to create output directory: %v", err))
	}

	outputData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON: %v", err))
	}

	err = os.WriteFile(config.OutputFile, outputData, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write output file: %v", err))
	}

	fmt.Println("Processing completed successfully! Output saved to:", config.OutputFile)
}
