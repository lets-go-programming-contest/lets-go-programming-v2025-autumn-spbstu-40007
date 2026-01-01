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
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
	"strconv"
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

	config := loadConfig(*configPath)
	valutes := loadXML(config.InputFile)
	result := transform(valutes)
	saveJSON(config.OutputFile, result)
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func loadConfig(path string) Config {
	data, err := os.ReadFile(path)
	check(err)

	var cfg Config
	check(yaml.Unmarshal(data, &cfg))

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		fmt.Fprintln(os.Stderr, "config fields must not be empty")
		os.Exit(1)
	}

	return cfg
}

func loadXML(path string) []Valute {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return input, nil
		}
	}

	var curs ValCurs
	check(decoder.Decode(&curs))

	return curs.Valutes
}

func transform(valutes []Valute) []ResultCurrency {
	result := make([]ResultCurrency, 0, len(valutes))

	for _, valute := range valutes {
		valueStr := strings.ReplaceAll(valute.Value, ",", ".")
		value, err := strconv.ParseFloat(valueStr, 64)
		check(err)

		result = append(result, ResultCurrency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CharCode < result[j].CharCode
	})

	return result
}

func saveJSON(path string, data []ResultCurrency) {
	dir := filepath.Dir(path)
	check(os.MkdirAll(dir, dirPerm))

	file, err := os.Create(path)
	check(err)
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	check(encoder.Encode(data))
}