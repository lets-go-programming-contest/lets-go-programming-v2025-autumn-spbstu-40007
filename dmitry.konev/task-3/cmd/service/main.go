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
	"gopkg.in/yaml.v3"
)

const dirPerm = 0o755

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type ResultCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "config flag is required")
		os.Exit(1)
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	valutes, err := loadXML(config.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load XML: %v\n", err)
		os.Exit(1)
	}

	result := transform(valutes)

	if err := saveJSON(config.OutputFile, result); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save JSON: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read file: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal yaml: %w", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return Config{}, fmt.Errorf("config fields must not be empty")
	}

	return cfg, nil
}

func loadXML(path string) ([]Valute, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}

	var curs ValCurs

	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return curs.Valutes, nil
}

func transform(valutes []Valute) []ResultCurrency {
	result := make([]ResultCurrency, 0, len(valutes))

	for _, valute := range valutes {
		if valute.CharCode == "" || valute.Value == "" {
			continue
		}

		value := parseFloat(valute.Value)

		result = append(result, ResultCurrency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Value > result[j].Value
	})

	return result
}

func saveJSON(path string, data []ResultCurrency) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("mkdir all: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

func parseFloat(s string) float64 {
	value, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0
	}
	return value
}