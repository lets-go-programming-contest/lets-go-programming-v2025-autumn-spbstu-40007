package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"

	"github.com/falsefeelings/task-3/iternal/config"
	"github.com/falsefeelings/task-3/iternal/data"
	"github.com/falsefeelings/task-3/iternal/pathmaker"
)

func main() {
	// 1. Конфиг
	cfg, err := config.Read("config.yaml")
	if err != nil {
		log.Fatal("Config error:", err)
	}

	// 2. Читаем XML
	xmlData, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		log.Fatal("Read file error:", err)
	}

	// 3. Парсим XML (используя data пакет)
	var valCourse data.ValCourse
	err = xml.Unmarshal(xmlData, &valCourse)
	if err != nil {
		log.Fatal("XML parse error:", err)
	}

	// 4. Конвертируем в JSON (встроенная функция)
	jsonData, err := json.MarshalIndent(valCourse, "", "  ")
	if err != nil {
		log.Fatal("JSON conversion error:", err)
	}

	// 5. Сохраняем
	pathmaker.CreateOutPath(cfg.OutputFile)
	err = os.WriteFile(cfg.OutputFile, jsonData, 0644)
	if err != nil {
		log.Fatal("Write file error:", err)
	}

	log.Println("Success! Converted", cfg.InputFile, "→", cfg.OutputFile)
}
