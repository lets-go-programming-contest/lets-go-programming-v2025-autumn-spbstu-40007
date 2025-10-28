package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/UwUshkin/task-3/internal/xmldecoder"
)

const outputPermissions = 0o600

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

	if err := os.WriteFile(outputPath, jsonData, outputPermissions); err != nil {
		return fmt.Errorf("writing output file %q: %w", outputPath, err)
	}

	return nil
}
