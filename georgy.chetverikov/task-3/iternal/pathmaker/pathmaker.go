package pathmaker

import (
	"fmt"
	"os"
	"path/filepath"
)

const permissions = 0o755

func CreateOutPath(outputFile string) error {
	directory := filepath.Dir(outputFile)

	err := os.MkdirAll(directory, permissions)
	if err != nil {
		return fmt.Errorf("creating directory for an output file at %q: %w", directory, err)
	}

	return nil
}
