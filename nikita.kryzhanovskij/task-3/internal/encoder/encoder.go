package encoder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"nikita.kryzhanovskij/task-3/internal/models"
)

func EncodeJSON(path string, data []models.ValuteOutput) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
