package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/UwUshkin/task-3/internal/xmldecoder"
)

func ProcessAndSave(inputPath, outputPath string) error {
	valCursData, err := xmldecoder.DecodeCBRXML(inputPath)
	if err != nil {
		return fmt.Errorf("decoding XML from %q: %w", inputPath, err)
	}

	sort.Sort(valCursData.Valutes)

	jsonData, err := json.MarshalIndent(valCursData.Valutes, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling results to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0o644); err != nil {
		return fmt.Errorf("writing output file %q: %w", outputPath, err)
	}

	return nil
}
