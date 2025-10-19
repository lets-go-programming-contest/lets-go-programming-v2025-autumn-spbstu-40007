package processor

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tntkatz/task-3/internal/config"
	"github.com/tntkatz/task-3/internal/crp"
	"github.com/tntkatz/task-3/internal/data"
	"github.com/tntkatz/task-3/internal/pathcreator"
	"github.com/tntkatz/task-3/internal/vp"
	"github.com/tntkatz/task-3/internal/xmldecoder"
)

const DefaultFilePermissions = 0o600

func Run(configPath string) error {
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	inputFile, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file '%s': %w", cfg.InputFile, err)
	}

	valCurs := data.ValCurs{
		Date:   "",
		Name:   "",
		Valute: nil,
	}

	err = xmldecoder.DecodeXML(inputFile, &valCurs)
	if err != nil {
		return fmt.Errorf("failed during XML decoding: %w", err)
	}

	sortedValutes, err := vp.ValuteProcess(valCurs)
	if err != nil {
		return fmt.Errorf("failed during valute processing: %w", err)
	}

	currencyResults := crp.CurrencyProcess(sortedValutes)

	jsonData, err := json.Marshal(currencyResults)
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	err = pathcreator.CreatePath(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output path: %w", err)
	}

	err = os.WriteFile(cfg.OutputFile, jsonData, DefaultFilePermissions)
	if err != nil {
		return fmt.Errorf("failed to write output file '%s': %w", cfg.OutputFile, err)
	}

	return nil
}
