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

	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load config:", err)
		os.Exit(1)
	}

	valutes, err := loadXML(cfg.InputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load XML:", err)
		os.Exit(1)
	}

	result := transform(valutes)

	if err := saveJSON(cfg.OutputFile, result); err != nil {
		fmt.Fprintln(os.Stderr, "save JSON:", err)
		os.Exit(1)
	}
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return Config{}, fmt.Errorf("config fields must not be empty")
	}

	return cfg, nil
}

func loadXML(path string) ([]Valute, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return curs.Valutes, nil
}

func transform(valutes []Valute) []ResultCurrency {
	result := make([]ResultCurrency, 0, len(valutes))

	for _, v := range valutes {
		if v.CharCode == "" {
			continue
		}
		result = append(result, ResultCurrency{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    parseFloat(v.Value),
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
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func parseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	if s == "" {
		return 0
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return value
}