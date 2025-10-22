package pathcreator

import (
	"fmt"
	"os"
	"path/filepath"
)

const DefaultDirPermissions = 0o755

func CreatePath(outputFile string) error {
	outputDir := filepath.Dir(outputFile)

	err := os.MkdirAll(outputDir, DefaultDirPermissions)
	if err != nil {
		return fmt.Errorf("creating output directory %q: %w", outputDir, err)
	}

	return nil
}
