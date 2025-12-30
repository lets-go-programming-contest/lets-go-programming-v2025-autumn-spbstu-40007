package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

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
		panic("config flag is required")
	}

	config := loadConfig(*configPath)
	valutes := loadXML(config.InputFile)
	result := transform(valutes)
	saveJSON(config.OutputFile, result)
}

func loadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		panic("config fields must not be empty")
	}

	return cfg
}

func loadXML(path string) []Valute {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var curs ValCurs
	if err := decoder.Decode(&curs); err != nil {
		panic(err)
	}

	return curs.Valutes
}

func transform(valutes []Valute) []ResultCurrency {
	result := make([]ResultCurrency, 0, len(valutes))

	for _, valute := range valutes {
		valueStr := strings.Replace(valute.Value, ",", ".", 1)
		value, err := parseFloat(valueStr)
		if err != nil {
			panic(err)
		}

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

func saveJSON(path string, data []ResultCurrency) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		panic(err)
	}
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscan(s, &f)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}

	return f, nil
}
