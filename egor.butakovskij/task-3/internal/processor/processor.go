package processor

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tntkatz/task-3/internal/config"
	"github.com/tntkatz/task-3/internal/data"
	"github.com/tntkatz/task-3/internal/pathcreator"
	"github.com/tntkatz/task-3/internal/vp"
	"github.com/tntkatz/task-3/internal/xmldecoder"
)

const DefaultFilePermissions = 0o600

func Run(configPath string) error {
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	inputFile, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("reading input file %q: %w", cfg.InputFile, err)
	}

	valCurs := data.ValCurs{} //nolint:exhaustivestruct

	err = xmldecoder.DecodeXML(inputFile, &valCurs)
	if err != nil {
		return fmt.Errorf("decoding XML: %w", err)
	}

	sortedValutes, err := vp.ValuteProcess(valCurs)
	if err != nil {
		return fmt.Errorf("processing valute: %w", err)
	}

	jsonData, err := json.Marshal(sortedValutes)
	if err != nil {
		return fmt.Errorf("marshalling results to JSON: %w", err)
	}

	err = pathcreator.CreatePath(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("creating output path: %w", err)
	}

	err = os.WriteFile(cfg.OutputFile, jsonData, DefaultFilePermissions)
	if err != nil {
		return fmt.Errorf("writing output file %q: %w", cfg.OutputFile, err)
	}

	return nil
}
